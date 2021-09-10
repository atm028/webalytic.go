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
	componentName := "handler"
	commonCfgOption := CommonCfg.Container()
	logger := fx.Option(fx.Provide(func(v *viper.Viper) bunyan.Logger {
		l := CommonCfg.Logger{
			Cfg: &CommonCfg.LoggerConfig{
				Viper: v,
			},
		}
		return l.GetLogger(componentName)
	}))
	handlerCfgOption := fx.Options(fx.Provide(func(
		v *viper.Viper,
	) *AppConfig.LogHandlerConfig {
		return &AppConfig.LogHandlerConfig{
			Viper: v,
		}
	}))

	redisBroker := RedisBroker.Container(componentName)
	clickhouse := Datasources.Clickhouse()
	ackRedisChannel := fx.Options(fx.Provide(func() chan string {
		ch := make(chan string)
		return ch
	}))

	return fx.Options(
		logger,
		commonCfgOption,
		handlerCfgOption,
		redisBroker,
		clickhouse,
		ackRedisChannel,
	)
}
