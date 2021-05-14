package RedisBroker

import (
	"fmt"
	"time"

	inversify "github.com/alekns/go-inversify"
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	CommonCfg "github.com/webalytic.go/common/config"
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
	redis   *redis.Client
	clients map[string]ClientInfo
	logger  bunyan.Logger
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

func Init(container inversify.Container) *RedisBroker {
	logger := CommonCfg.GetLogger(container)
	logger.Debug("Init RedisStreamEventBroker")
	cfgObj, _ := container.Get("redisConfig")
	cfg, _ := cfgObj.(*CommonCfg.RedisConfig)
	logger.Debug("Redis broker: ", cfg)

	broker := &RedisBroker{
		redis: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", cfg.Host(), cfg.Port()),
		}),
		clients: make(map[string]ClientInfo),
		logger:  logger,
	}
	logger.Debug("Created broker on redis server: ", broker.redis.Options().Addr)
	go broker.readStream()

	return broker
}
