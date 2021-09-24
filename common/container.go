package app

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	CommonCfg "github.com/webalytic.go/common/config"
	Consul "github.com/webalytic.go/common/consul"
	"go.uber.org/fx"
)

func Container(name string) fx.Option {
	logger := fx.Option(fx.Provide(func(v *viper.Viper) bunyan.Logger {
		l := CommonCfg.Logger{
			Cfg: &CommonCfg.LoggerConfig{
				Viper: v,
			},
		}
		return l.GetLogger(name)
	}))

	consul := Consul.Container()

	return fx.Options(logger, consul)
}
