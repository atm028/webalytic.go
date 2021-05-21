package app

import (
	"fmt"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/go-redis/redis"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func RedisEventBrokerHandler(
	logger bunyan.Logger,
	broker *RedisBroker.RedisBroker,
	cfg *AppConfig.LogHandlerConfig,
	evtChannel chan redis.XMessage) {
	for {
		event := <-evtChannel
		fmt.Println(event)
		logger.Debug(event.Values)
		broker.Ack(cfg.Channel(), event.ID)
	}
}
