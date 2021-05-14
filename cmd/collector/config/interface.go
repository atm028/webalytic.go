package config

import (
	Inversify "github.com/alekns/go-inversify"
	CommonCfg "github.com/webalytic.go/common/config"
)

type ICollectorConfig interface {
	Port() int
	Level() string
	Channel() string
}

type CollectorConfig struct {
	Container Inversify.Container
}

func (c CollectorConfig) Port() int {
	viper := CommonCfg.GetViper(c.Container)

	viper.SetDefault("collector.port", 8090)
	port := viper.GetInt("collector.port")
	if port == 0 {
		viper.BindEnv("port")
		port = viper.GetInt("port")
	}
	return port
}

func (c CollectorConfig) Level() string {
	viper := CommonCfg.GetViper(c.Container)
	viper.SetDefault("collector.level", "DEBUG")
	level := viper.GetString("collector.level")
	if len(level) == 0 {
		viper.BindEnv("level")
		level = viper.GetString("level")
	}
	return level
}

func (c CollectorConfig) Channel() string {
	viper := CommonCfg.GetViper(c.Container)
	viper.SetDefault("collector.channel", "COLLECTOR")
	ch := viper.GetString("collector.channel")
	if len(ch) == 0 {
		viper.BindEnv("channel")
		ch = viper.GetString("channel")
	}
	return ch
}
