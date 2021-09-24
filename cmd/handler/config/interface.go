package config

import (
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/spf13/viper"
)

type IHandlerConfig interface {
	Name() string
	Port() int
	Level() string
	Channel() string
}

type HandlerConfig struct {
	Viper  *viper.Viper
	logger bunyan.Logger
}

func (c HandlerConfig) Name() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("name", "handler")
	c.Viper.BindEnv("name")
	name := c.Viper.GetString("name")
	if len(name) == 0 {
		name = c.Viper.GetString("handler.name")
	}
	return name
}

func (c HandlerConfig) Port() int {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("port", 8091)
	port := c.Viper.GetInt("port")
	if port == 0 {
		c.Viper.BindEnv("port")
		port = c.Viper.GetInt("handler.port")
	}
	return port
}

func (c HandlerConfig) Level() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("level", "DEBUG")
	level := c.Viper.GetString("level")
	if len(level) == 0 {
		c.Viper.BindEnv("level")
		level = c.Viper.GetString("handler.level")
	}
	return level
}

func (c HandlerConfig) Channel() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("channel", "COLLECTOR")
	ch := c.Viper.GetString("channel")
	if len(ch) == 0 {
		c.Viper.BindEnv("channel")
		ch = c.Viper.GetString("handler.channel")
	}
	return ch
}
