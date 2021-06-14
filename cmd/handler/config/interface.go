package config

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
)

type ILogHandlerConfig interface {
	Name() string
	Port() int
	Level() string
	Channel() string
}

type LogHandlerConfig struct {
	Viper  *viper.Viper
	logger bunyan.Logger
}

func (c LogHandlerConfig) Name() string {
	c.Viper.SetDefault("handler.name", "loghandler")
	c.Viper.BindEnv("name")
	name := c.Viper.GetString("name")
	if len(name) == 0 {
		name = c.Viper.GetString("handler.name")
	}
	return name
}

func (c LogHandlerConfig) Port() int {
	c.Viper.SetDefault("handler.port", 8091)
	port := c.Viper.GetInt("handler.port")
	if port == 0 {
		c.Viper.BindEnv("port")
		port = c.Viper.GetInt("port")
	}
	return port
}

func (c LogHandlerConfig) Level() string {
	c.Viper.SetDefault("handler.level", "DEBUG")
	level := c.Viper.GetString("handler.level")
	if len(level) == 0 {
		c.Viper.BindEnv("level")
		level = c.Viper.GetString("level")
	}
	return level
}

func (c LogHandlerConfig) Channel() string {
	c.Viper.SetDefault("handler.channel", "COLLECTOR")
	ch := c.Viper.GetString("handler.channel")
	if len(ch) == 0 {
		c.Viper.BindEnv("channel")
		ch = c.Viper.GetString("channel")
	}
	return ch
}
