package app

import (
	"fmt"

	Inversify "github.com/alekns/go-inversify"
	"github.com/go-redis/redis"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func RedisEventBrokerHandler(
	container Inversify.Container,
	evtChannel chan redis.XMessage) {
	cfgObj, _ := container.Get("appConfig")
	cfg, _ := cfgObj.(*AppConfig.LogHandlerConfig)
	cfgObj, _ = container.Get("redisBroker")
	broker, _ := cfgObj.(*RedisBroker.RedisBroker)
	logger := CommonCfg.GetLogger(container)

	for {
		event := <-evtChannel
		fmt.Println(event)
		logger.Debug(event.Values)
		broker.Ack(cfg.Channel(), event.ID)
	}
}
