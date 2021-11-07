package config

import (
	"github.com/spf13/viper"
)

type IRedisConfig interface {
	Port() int
	Host() string
	StreamName() string
	HandlerName() string
	GroupName() string
}
type RedisConfig struct {
	Viper *viper.Viper
}

func (c RedisConfig) Port() int {
	c.Viper.SetEnvPrefix("REDIS")
	c.Viper.BindEnv("port")
	port := c.Viper.GetInt("port")
	if port == 0 {
		c.Viper.SetDefault("redis.port", 6379)
		port = c.Viper.GetInt("redis.port")
	}
	return port
}

func (c RedisConfig) Host() string {
	c.Viper.SetEnvPrefix("REDIS")
	c.Viper.BindEnv("host")
	host := c.Viper.GetString("host")

	if len(host) == 0 {
		c.Viper.SetDefault("redis.host", "0.0.0.0")
		host = c.Viper.GetString("redis.host")
	}
	return host
}

func (c RedisConfig) HandlerName() string {
	c.Viper.SetEnvPrefix("REDIS")
	c.Viper.BindEnv("handler")
	name := c.Viper.GetString("handler")

	if len(name) == 0 {
		c.Viper.SetDefault("redis.handler", "")
		name = c.Viper.GetString("redis.handler")
	}
	return name
}

func (c RedisConfig) StreamName() string {
	c.Viper.SetEnvPrefix("REDIS")
	c.Viper.BindEnv("stream")
	name := c.Viper.GetString("stream")

	if len(name) == 0 {
		c.Viper.SetDefault("redis.stream", "collector-stream")
		name = c.Viper.GetString("redis.stream")
	}
	return name
}

func (c RedisConfig) GroupName() string {
	c.Viper.SetEnvPrefix("REDIS")
	c.Viper.BindEnv("group")
	name := c.Viper.GetString("group")

	if len(name) == 0 {
		c.Viper.SetDefault("redis.group", "")
		name = c.Viper.GetString("redis.group")
	}
	return name
}
