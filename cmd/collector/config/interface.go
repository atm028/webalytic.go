package config

import (
	"github.com/spf13/viper"
)

type ICollectorConfig interface {
	Name() string
	Port() int
	Level() string
	Channel() string
}

type CollectorConfig struct {
	Viper *viper.Viper
}

func (c *CollectorConfig) Name() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("name", "collector")
	name := c.Viper.GetString("name")
	return name
}

func (c *CollectorConfig) Port() int {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("port", 8090)
	port := c.Viper.GetInt("port")
	return port
}

func (c *CollectorConfig) Level() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("level", "DEBUG")
	level := c.Viper.GetString("level")
	if len(level) == 0 {
		c.Viper.BindEnv("level")
		level = c.Viper.GetString("collector.level")
	}
	return level
}

func (c *CollectorConfig) Channel() string {
	c.Viper.SetEnvPrefix("APP")
	c.Viper.SetDefault("channel", "COLLECTOR")
	ch := c.Viper.GetString("channel")
	if len(ch) == 0 {
		c.Viper.BindEnv("channel")
		ch = c.Viper.GetString("collector.channel")
	}
	return ch
}
