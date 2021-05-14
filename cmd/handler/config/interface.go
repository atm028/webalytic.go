package config

import (
	Inversify "github.com/alekns/go-inversify"
	CommonCfg "github.com/webalytic.go/common/config"
)

type ILogHandlerConfig interface {
	Port() int
	Level() string
	Channel() string
}

type LogHandlerConfig struct {
	Container Inversify.Container
}

func (c LogHandlerConfig) Port() int {
	viper := CommonCfg.GetViper(c.Container)
	viper.SetDefault("handler.port", 8091)
	port := viper.GetInt("handler.port")
	if port == 0 {
		viper.BindEnv("port")
		port = viper.GetInt("port")
	}
	return port
}

func (c LogHandlerConfig) Level() string {
	viper := CommonCfg.GetViper(c.Container)
	viper.SetDefault("handler.level", "DEBUG")
	level := viper.GetString("handler.level")
	if len(level) == 0 {
		viper.BindEnv("level")
		level = viper.GetString("level")
	}
	return level
}

func (c LogHandlerConfig) Channel() string {
	viper := CommonCfg.GetViper(c.Container)
	viper.SetDefault("handler.channel", "COLLECTOR")
	ch := viper.GetString("handler.channel")
	if len(ch) == 0 {
		viper.BindEnv("channel")
		ch = viper.GetString("channel")
	}
	return ch
}
