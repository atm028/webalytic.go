package app

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/collector/config"
	CommonCfg "github.com/webalytic.go/common/config"
	Consul "github.com/webalytic.go/common/consul"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func Container() fx.Option {
	name := fx.Option(fx.Provide(func(cfg *AppConfig.CollectorConfig) string {
		return cfg.Name()
	}))
	commonCfgOption := CommonCfg.Container()
	consul := Consul.Container()

	logger := fx.Option(fx.Provide(func(
		v *viper.Viper,
		cfg *AppConfig.CollectorConfig) bunyan.Logger {
		l := CommonCfg.Logger{
			Cfg: &CommonCfg.LoggerConfig{
				Viper: v,
			},
		}
		return l.GetLogger(cfg.Name())
	}))

	InitCollectHandler()
	handlerOption := fx.Options(fx.Provide(CollectHandler))
	healthOption := fx.Option(fx.Provide(HealthHandler))

	collectorCfgOption := fx.Options(fx.Provide(func(v *viper.Viper) *AppConfig.CollectorConfig {
		return &AppConfig.CollectorConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container()

	return fx.Options(
		logger,
		name,
		consul,
		commonCfgOption,
		handlerOption,
		healthOption,
		collectorCfgOption,
		redisBroker,
	)
}
