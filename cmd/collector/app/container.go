package app

import (
	Inversify "github.com/alekns/go-inversify"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
)

func Container() Inversify.Container {
	container := Inversify.NewContainer("collector")
	container.Bind("name").To("collector")
	container = CommonCfg.Container(container)
	container.Bind("collect").To(CollectHandler).InSingletonScope()
	container.Bind("appConfig").To(&AppConfig.CollectorConfig{Container: container})
	container.Bind("redisConfig").To(&CommonCfg.RedisConfig{Container: container})
	container = RedisBroker.Container(container)
	return container
}
