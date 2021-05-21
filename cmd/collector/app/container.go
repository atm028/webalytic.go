package app

import (
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func Container() fx.Option {
	appCfgOption := fx.Options(fx.Provide(func() *CommonCfg.AppConfig {
		return &CommonCfg.AppConfig{
			Name: "collector",
		}
	}))
	commonCfgOption := CommonCfg.Container()
	handlerOption := fx.Options(fx.Provide(CollectHandler))
	collectorCfgOption := fx.Options(fx.Provide(func(v *viper.Viper) *AppConfig.CollectorConfig {
		return &AppConfig.CollectorConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container()

	return fx.Options(
		appCfgOption,
		commonCfgOption,
		handlerOption,
		collectorCfgOption,
		redisBroker,
	)
}
