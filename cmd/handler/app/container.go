package app

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	AppConfig "github.com/webalytic.go/cmd/handler/config"
	CommonCfg "github.com/webalytic.go/common/config"
	Datasources "github.com/webalytic.go/common/datasources"
	RedisBroker "github.com/webalytic.go/common/redis"
	"go.uber.org/fx"
)

func Container() fx.Option {
	name := fx.Option(fx.Provide(func(cfg *AppConfig.HandlerConfig) string {
		return cfg.Name()
	}))
	commonCfgOption := CommonCfg.Container()
	logger := fx.Option(fx.Provide(func(
		v *viper.Viper,
		cfg *AppConfig.HandlerConfig) bunyan.Logger {
		l := CommonCfg.Logger{
			Cfg: &CommonCfg.LoggerConfig{
				Viper: v,
			},
		}
		return l.GetLogger(cfg.Name())
	}))
	handlerCfgOption := fx.Options(fx.Provide(func(
		v *viper.Viper,
	) *AppConfig.HandlerConfig {
		return &AppConfig.HandlerConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container()
	clickhouse := Datasources.Clickhouse()
	ackRedisChannel := fx.Options(fx.Provide(func() chan string {
		ch := make(chan string)
		return ch
	}))

	return fx.Options(
		logger,
		name,
		commonCfgOption,
		handlerCfgOption,
		redisBroker,
		clickhouse,
		ackRedisChannel,
	)
}
