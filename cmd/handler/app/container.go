package app

import (
	Inversify "github.com/alekns/go-inversify"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func Container() Inversify.Container {
	container := Inversify.NewContainer("handler")
	container.Bind("name").To("loghandler")
	container = CommonCfg.Container(container)
	container.Bind("appConfig").To(&AppConfig.LogHandlerConfig{Container: container})
	container.Bind("redisConfig").To(&CommonCfg.RedisConfig{Container: container})
	container = RedisBroker.Container(container)
	return container
}
