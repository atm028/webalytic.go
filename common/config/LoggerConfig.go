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
	name := c.Viper.GetString("name")
	fmt.Println("LoggerConfig: name: ", name)
	return name
}

func (c *LoggerConfig) Level() string {
	return "DEBUG"
}

func (c *LoggerConfig) Path() string {
	return "./"
}
