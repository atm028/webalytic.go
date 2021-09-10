package Datasources

//TODO: add migration : curl -XPOST http://0.0.0.0:8123/?query -d "create database if not exists webalytic"

import (
	"fmt"

	"sync"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type IClickHouse interface {
	CreatePayment(data *Payment) error
	FindPayment(req string, val string) (*Payment, error)
	flushTimer() *chan struct{}
	flush()
}

type ClickHouse struct {
	Db            *gorm.DB
	payments      []*Payment
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

	if c.count > 0 {
		res := c.Db.Create(&c.payments)
		if res.Error != nil {
			c.logger.Error(fmt.Sprintf("Payments insert error: %s", res.Error))
		} else {
			c.logger.Debug(fmt.Sprintf("Flushed %d records", c.count))
			for k := range c.payments {
				key := c.payments[k].CacheKey
				c.logger.Debug(fmt.Sprintf("Send ACK: %s", key))
				c.ackCh <- key
			}
			c.payments = c.payments[:0]
			c.count = 0
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

func (c *ClickHouse) CreatePayment(data *Payment) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payments = append(c.payments, data)
	c.count++
	if c.count >= c.FlushLimit {
		c.flush()
	}
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
		logger.Debug("Migration: Payment")
		db.AutoMigrate(&Payment{})
		db.Migrator().CreateTable(Payment{})
		logger.Debug("Migration: Payment: Done")
		ch := &ClickHouse{
			Db:            db,
			logger:        logger,
			payments:      nil,
			count:         0,
			ackCh:         ackRedisChannel,
			flushInterval: config.FlushInterval(),
			FlushLimit:    config.FlushLimit(),
		}
		go ch.flushTimer()
		return ch
	}))
}
