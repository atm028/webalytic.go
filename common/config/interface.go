package config

import (
	"fmt"
	"os"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AppConfig struct {
	Name string
}

type IRedisConfig interface {
	Port() int
	Host() string
}
type RedisConfig struct {
	Prefix string
	Viper  *viper.Viper
}

func (c RedisConfig) Port() int {
	c.Viper.SetEnvPrefix(c.Prefix)
	c.Viper.BindEnv("redis_port")
	port := c.Viper.GetInt("redis_port")
	if port == 0 {
		c.Viper.SetDefault("redis.port", 6379)
		port = c.Viper.GetInt("redis.port")
	}
	return port
}

func (c RedisConfig) Host() string {
	c.Viper.SetEnvPrefix(c.Prefix)
	c.Viper.BindEnv("redis_host")
	host := c.Viper.GetString("redis_host")

	if len(host) == 0 {
		c.Viper.SetDefault("redis.host", "0.0.0.0")
		host = c.Viper.GetString("redis.host")
	}
	return host
}

func Container() fx.Option {
	Logger := fx.Options(fx.Provide(func(cfg *AppConfig) bunyan.Logger {
		staticFields := make(map[string]interface{})
		logName := "./logs/" + cfg.Name + ".log"
		bunyanConfig := bunyan.Config{
			Name: cfg.Name,
			Streams: []bunyan.Stream{
				{
					Name:   cfg.Name,
					Level:  bunyan.LogLevelDebug,
					Stream: os.Stdout,
				},
				{
					Name:  cfg.Name,
					Level: bunyan.LogLevelDebug,
					Path:  logName,
				},
			},
			StaticFields: staticFields,
		}
		logger, err := bunyan.CreateLogger(bunyanConfig)
		if err != nil {
			fmt.Println("Error: cannot create logger")
		}
		logger.Info("set env prefix: ", cfg.Name)
		logger.Debug("Logger at " + logName + " initialized")

		return logger
	}))

	Viper := fx.Options(fx.Provide(func(logger bunyan.Logger) *viper.Viper {
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

		logger.Debug("Init common config")
		return v
	}))

	Redis := fx.Options(fx.Provide(func(cfg *AppConfig, v *viper.Viper) *RedisConfig {
		return &RedisConfig{
			Prefix: cfg.Name,
			Viper:  v,
		}
	}))

	return fx.Options(Logger, Viper, Redis)
}
