package Datasources

//TODO: add migration : curl -XPOST http://0.0.0.0:8123/?query -d "create database if not exists webalytic"

import (
	"fmt"
	"sync"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/prometheus/client_golang/prometheus"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var (
	handlerFlushedRecordsCnt = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "handler_flushed_records",
		Subsystem: "Handler",
		Help:      "Number of records flushed into the DB",
	})
)

type IClickHouse interface {
	CreatePayment(key string, payment *Payment) error
	FindPayment(req string, val string) (*Payment, error)
	flushTimer() *chan bool
	flush()
}

type ClickHouse struct {
	Db            *gorm.DB
	payments      map[string]*Payment
	count         int
	ackCh         chan string
	mu            sync.Mutex
	flushInterval int
	FlushLimit    int
	logger        bunyan.Logger
}

func (c *ClickHouse) flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Db == nil {
		c.logger.Fatal("DB is null")
	}

	if c.Db != nil && c.count > 0 {
		payments := make([]*Payment, 0, c.count)
		for _, v := range c.payments {
			payments = append(payments, v)
		}
		res := c.Db.Create(&payments)
		if res.Error != nil {
			c.logger.Error(fmt.Sprintf("Not able to create record: %s", res.Error))
		} else {
			c.logger.Debug(fmt.Sprintf("Created %d records", len(payments)))
			for key := range c.payments {
				c.logger.Debug(fmt.Sprintf("Send ACK for key: %s", key))
				c.ackCh <- key
				handlerFlushedRecordsCnt.Inc()
				delete(c.payments, key)
				c.count--
			}
		}
	}
}

func (c *ClickHouse) flushTimer() *chan bool {
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(c.flushInterval))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.flush()
			case <-done:
				return
			}
		}
	}()
	return &done
}

func (c *ClickHouse) CreatePayment(key string, payment *Payment) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payments[key] = payment
	c.count++
	if c.count >= c.FlushLimit {
		c.flush()
	}
	return nil
}

func (c *ClickHouse) FindPayment(req string, val string) (*Payment, error) {
	var res Payment
	c.Db.Find(&res, req, val)
	return &res, nil
}

func Clickhouse() fx.Option {
	return fx.Options(fx.Provide(func(
		logger bunyan.Logger,
		config *CommonCfg.ClickhouseConfig,
		ackRedisChannel chan string) *ClickHouse {
		logger.Debug(fmt.Sprintf("Connecting to CH: %s", config.NativeConnStr()))
		db, err := gorm.Open(clickhouse.Open(config.NativeConnStr()), &gorm.Config{})
		if err != nil {
			logger.Fatal(fmt.Sprintf("Not able to connect to Clickhouse: %s", err))
		}
		prometheus.MustRegister(handlerFlushedRecordsCnt)

		/**
		Provide migration for each suported type
		*/
		// 1. Payment
		logger.Debug("Migration: Payment")
		db.AutoMigrate(&Payment{})
		db.Migrator().CreateTable(Payment{})
		logger.Debug("Migration: Payment: Done")

		ch := &ClickHouse{
			Db:            db,
			logger:        logger,
			payments:      make(map[string]*Payment),
			count:         0,
			ackCh:         ackRedisChannel,
			flushInterval: config.FlushInterval(),
			FlushLimit:    config.FlushLimit(),
		}
		go ch.flushTimer()
		return ch
	}))
}
