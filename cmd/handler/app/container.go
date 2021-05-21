package app

import (
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func Container() fx.Option {
	appCfgOption := fx.Options(fx.Provide(func() *CommonCfg.AppConfig {
		return &CommonCfg.AppConfig{
			Name: "loghandler",
		}
	}))
	commonCfgOption := CommonCfg.Container()
	handlerCfgOption := fx.Options(fx.Provide((func(v *viper.Viper) *AppConfig.LogHandlerConfig {
		return &AppConfig.LogHandlerConfig{
			Viper: v,
		}
	})))

	redisBroker := RedisBroker.Container()

	return fx.Options(
		appCfgOption,
		commonCfgOption,
		handlerCfgOption,
		redisBroker,
	)
}
