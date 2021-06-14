package app

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
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
	handlerOption := fx.Options(fx.Provide(CollectHandler))
	collectorCfgOption := fx.Options(fx.Provide(func(v *viper.Viper) *AppConfig.CollectorConfig {
		return &AppConfig.CollectorConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container("collector")

	return fx.Options(
		logger,
		commonCfgOption,
		handlerOption,
		collectorCfgOption,
		redisBroker,
	)
}
