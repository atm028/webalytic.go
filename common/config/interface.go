package config

import (
	"fmt"
	"os"

	Inversify "github.com/alekns/go-inversify"
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
)

type IRedisConfig interface {
	Port() int
	Host() string
}
type RedisConfig struct {
	Container Inversify.Container
}

func (c RedisConfig) Port() int {
	v := GetViper(c.Container)
	cfgObj, _ := c.Container.Get("name")
	appName, _ := cfgObj.(string)
	v.SetEnvPrefix(appName)
	v.BindEnv("redis_port")
	port := v.GetInt("redis_port")
	if port == 0 {
		v.SetDefault("redis.port", 6379)
		port = v.GetInt("redis.port")
	}
	return port
}

func (c RedisConfig) Host() string {
	v := GetViper(c.Container)
	cfgObj, _ := c.Container.Get("name")
	appName, _ := cfgObj.(string)
	v.SetEnvPrefix(appName)
	v.BindEnv("redis_host")
	host := v.GetString("redis_host")

	if len(host) == 0 {
		v.SetDefault("redis.host", "0.0.0.0")
		host = v.GetString("redis.host")
	}
	return host
}

func GetLogger(container Inversify.Container) bunyan.Logger {
	cfgObj, err := container.Get("logger")
	if err != nil {
		fmt.Println("Error: no Cfg object fro logger")
	}
	logger, _ := cfgObj.(bunyan.Logger)
	return logger
}

func GetViper(container Inversify.Container) *viper.Viper {
	logger := GetLogger(container)
	cfgObj, err := container.Get("viper")
	if err != nil {
		logger.Fatal("Cannot get viper from container")
	}
	v := cfgObj.(*viper.Viper)
	return v
}

func Container(container Inversify.Container) Inversify.Container {
	v := viper.GetViper()
	v.SetConfigName("default")
	v.SetConfigType("yaml")
	v.AddConfigPath("~/go/src/github.com/webalytic.go/common/config")
	v.AddConfigPath("./")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("Unable to read config")
	}

	cfgObj, _ := container.Get("name")
	appName, _ := cfgObj.(string)
	fmt.Println("set env prefix: ", appName)
	staticFields := make(map[string]interface{})
	logName := "./logs/" + appName + ".log"
	bunyanConfig := bunyan.Config{
		Name: appName,
		Streams: []bunyan.Stream{
			{
				Name:   appName,
				Level:  bunyan.LogLevelDebug,
				Stream: os.Stdout,
			},
			{
				Name:  appName,
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
	logger.Debug("Logger at " + logName + " initialized")
	container.Bind("logger").To(logger).InSingletonScope()

	v.SetDefault("APP_PREFIX", "app")
	v.SetEnvPrefix(v.GetString("APP_PREFIX"))
	container.Bind("viper").To(v).InSingletonScope()
	logger.Debug("Init common config")
	return container
}
