package app

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func Container() fx.Option {
	commonCfgOption := CommonCfg.Container()
	logger := fx.Option(fx.Provide(func(v *viper.Viper) bunyan.Logger {
		l := CommonCfg.Logger{
			Cfg: &CommonCfg.LoggerConfig{
				Viper: v,
			},
		}
		return l.GetLogger("collector")
	}))
	handlerCfgOption := fx.Options(fx.Provide(func(
		v *viper.Viper,
	) *AppConfig.LogHandlerConfig {
		return &AppConfig.LogHandlerConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container("handler")

	return fx.Options(
		logger,
		commonCfgOption,
		handlerCfgOption,
		redisBroker,
	)
}
