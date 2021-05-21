package RedisBroker

import (
	"fmt"
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"
)

type IRedisEventBroker interface {
	Subscribe(channel string, evtChannel chan redis.XMessage)
	Unsubscribe(channel string)
	Publish(channel string, msg []byte)
	Ack(channel string, id string)
}

type ClientInfo struct {
	evtChannel chan redis.XMessage
	lastID     string
}

type RedisBroker struct {
	redis    *redis.Client
	clients  map[string]ClientInfo
	logger   bunyan.Logger
	consumer string
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

func (broker *RedisBroker) readStream() {
	for {
		StreamsTmplt := []string{}
		for k, v := range broker.clients {
			StreamsTmplt = append(StreamsTmplt, k, v.lastID)
		}

		entries, _ := broker.redis.XRead(&redis.XReadArgs{
			Streams: StreamsTmplt,
			Count:   1,
			Block:   100 * time.Millisecond,
		}).Result()

		for _, entry := range entries {
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

	for {
		entries, _ := broker.redis.XReadGroup(&redis.XReadGroupArgs{
			Group:    "group",
			Consumer: broker.consumer,
			Streams:  []string{"stream", ">"},
		}).Result()
		for _, entry := range entries {
			broker.logger.Info("readGroupStream: event: ", entry)
		}
	}
}

func (broker *RedisBroker) createGroupService(channel string, group string) error {
	err := broker.redis.XGroupCreate(channel, group, "0").Err()
	if err != nil {
		broker.logger.Error("Error: createGroupService: ", err)
	}
	return err
}

func Container() fx.Option {
	return fx.Options(fx.Provide(func(
		logger bunyan.Logger,
		commonAppCfg *CommonCfg.AppConfig,
		redisCfg *CommonCfg.RedisConfig,
	) *RedisBroker {
		logger.Debug("Init RedisStreamEventBroker")
		logger.Debug("Redis broker: ", redisCfg)

		broker := &RedisBroker{
			redis: redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", redisCfg.Host(), redisCfg.Port()),
			}),
			clients:  make(map[string]ClientInfo),
			logger:   logger,
			consumer: commonAppCfg.Name,
		}
		logger.Debug("Created broker on redis server: ", broker.redis.Options().Addr)
		go broker.readStream()
		go broker.readGroupStream()

		return broker
	}))
}
