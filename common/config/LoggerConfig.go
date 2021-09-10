package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ILoggerConfig interface {
	Name() string
	Level() string
	Path() string
}

type LoggerConfig struct {
	Viper *viper.Viper
}

func (c *LoggerConfig) Name() string {
	c.Viper.BindEnv("name")
	c.Viper.SetDefault("name", "webalytic")
	name := c.Viper.GetString("name")
	fmt.Println("LoggerConfig: name: ", name)
	return name
}

func (c *LoggerConfig) Level() string {
	c.Viper.SetEnvPrefix("LOG")
	c.Viper.BindEnv("level")
	c.Viper.SetDefault("level", "DEBUG")
	level := c.Viper.GetString("level")
	return level
}

func (c *LoggerConfig) Path() string {
	c.Viper.SetEnvPrefix("LOG")
	c.Viper.BindEnv("path")
	c.Viper.SetDefault("path", "./")
	v := c.Viper.GetString("path")
	return v
}
