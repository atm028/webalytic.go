package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Container() fx.Option {
	Viper := fx.Options(fx.Provide(func() *viper.Viper {
		v := viper.GetViper()
		v.SetConfigName("default")
		v.SetConfigType("yaml")
		v.AddConfigPath("./common/config")
		v.AddConfigPath("./")
		err := v.ReadInConfig()
		if err != nil {
			fmt.Println("Unable to read config")
		}
		v.SetDefault("APP_PREFIX", "app")
		v.SetEnvPrefix(v.GetString("APP_PREFIX"))

		fmt.Println("Init common config")
		return v
	}))

	Redis := fx.Options(fx.Provide(func(v *viper.Viper) *RedisConfig {
		return &RedisConfig{
			Viper: v,
		}
	}))

	LoggerObj := fx.Options(fx.Provide(func(v *viper.Viper) *Logger {
		return &Logger{
			Cfg: &LoggerConfig{
				Viper: v,
			},
		}
	}))

	return fx.Options(Viper, Redis, LoggerObj)
}
