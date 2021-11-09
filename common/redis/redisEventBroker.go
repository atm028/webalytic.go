package RedisBroker

import (
	"context"
	"fmt"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"
)

var (
	handlerMessagesNumberFromRedis = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "handler_number_messages_from_redis",
		Subsystem: "Handler",
		Help:      "Number of messages received from Redis by Handler",
	})
)

type IRedisEventBroker interface {
	Subscribe(channel string, evtChannel chan redis.XMessage)
	Unsubscribe(channel string)
	Publish(channel string, msg []byte)
	Ack(channel string, id string)
	GroupAck(id string)
}

type ClientInfo struct {
	evtChannel chan redis.XMessage
	lastID     string
}

type RedisBroker struct {
	redis       *redis.Client
	clients     map[string]ClientInfo
	logger      bunyan.Logger
	consumer    string
	streamName  string
	handlerName string
}

func (broker *RedisBroker) Subscribe(channel string, evtChannel chan redis.XMessage) {
	broker.clients[channel] = ClientInfo{evtChannel: evtChannel, lastID: "0"}
	broker.logger.Debug("Subscribed: ", channel)
}

func (broker *RedisBroker) Unsubscribe(channel string) {
	delete(broker.clients, channel)
}

func (broker *RedisBroker) Publish(channel string, msg []byte) error {
	broker.logger.Debug("publish: channel: ", channel)
	err := broker.redis.XAdd(&redis.XAddArgs{
		Stream:       channel,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"msg": string(msg),
		},
	}).Err()
	if err != nil {
		broker.logger.Error("Error: ", err)
	}
	return err
}

func (broker *RedisBroker) Ack(channel string, id string) {
	broker.redis.XDel(channel, id)
}

func (broker *RedisBroker) GroupAck(id string) {
	broker.logger.Debug(fmt.Sprintf("ACK for %s", id))
	broker.redis.XAck(broker.streamName, broker.handlerName, id)
}

func (broker *RedisBroker) readStream() {
	broker.logger.Info(fmt.Sprintf("readStream:starting, streamName: %s", broker.streamName))
	broker.logger.Debug("readStream:clients:", broker.clients)
	for {
		StreamsTmplt := []string{}
		for k, v := range broker.clients {
			if k != broker.streamName {
				StreamsTmplt = append(StreamsTmplt, k, v.lastID)
			}
		}

		entries, _ := broker.redis.XRead(&redis.XReadArgs{
			Streams: StreamsTmplt,
			Count:   1,
			Block:   100 * time.Millisecond,
		}).Result()

		for _, entry := range entries {
			broker.logger.Debug("readStream: entry: ", entry)
			client := broker.clients[entry.Stream]
			messages := entry.Messages
			for _, msg := range messages {
				client.lastID = msg.ID
				broker.clients[entry.Stream] = client
				client.evtChannel <- msg
			}
		}
	}
}

func (broker *RedisBroker) readGroupStream() {
	broker.logger.Info("readGroupStream:starting for consumer: %s, group: %s, streams: %s", broker.consumer, "log-handlers", "collector-stream")
	for {
		entries, _ := broker.redis.XReadGroup(&redis.XReadGroupArgs{
			Group:    "log-handlers",
			Consumer: broker.consumer,
			Streams:  []string{broker.streamName, ">"},
		}).Result()
		for _, entry := range entries {
			broker.logger.Debug("readGroupStream: entry: ", entry)
			client := broker.clients[entry.Stream]
			messages := entry.Messages
			handlerMessagesNumberFromRedis.Set(float64(len(messages)))
			for _, msg := range messages {
				broker.logger.Debug("send to channel")
				client.evtChannel <- msg
			}
		}
	}
}

func (broker *RedisBroker) createGroupStream(channel string, group string) {
	stat := broker.redis.XGroupCreateMkStream(channel, group, "$").Err()
	broker.logger.Debug("createGroupService: stat: ", stat)
}

func Container() fx.Option {
	return fx.Options(fx.Provide(func(
		lc fx.Lifecycle,
		logger bunyan.Logger,
		redisCfg *CommonCfg.RedisConfig,
		name string) *RedisBroker {
		logger.Debug("Init RedisStreamEventBroker with consumer name = %s", name)
		logger.Debug("Redis broker: ", redisCfg)

		lc.Append(
			fx.Hook{
				OnStop: func(ctx context.Context) error {
					logger.Debug("RedisEventBroker: OnStop hook")
					prometheus.Unregister(handlerMessagesNumberFromRedis)

					return nil
				},
			},
		)

		logger.Debug(fmt.Sprintf("Create RedisBroker wih streamName = %s on server %s:%d",
			redisCfg.StreamName(),
			redisCfg.Host(),
			redisCfg.Port()))
		broker := &RedisBroker{
			redis: redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", redisCfg.Host(), redisCfg.Port()),
			}),
			clients:     make(map[string]ClientInfo),
			logger:      logger,
			consumer:    name,
			streamName:  redisCfg.StreamName(),
			handlerName: redisCfg.HandlerName(),
		}
		logger.Debug("Created broker on redis server: ", broker.redis.Options().Addr)

		broker.createGroupStream(redisCfg.StreamName(), redisCfg.GroupName())
		logger.Debug("createGroupStream done")

		prometheus.MustRegister(handlerMessagesNumberFromRedis)

		/*
			TODO: publication to stream and group at the same time,
			so duplicates records in Clickhouse
		*/
		go broker.readStream()
		go broker.readGroupStream()

		return broker
	}))
}
