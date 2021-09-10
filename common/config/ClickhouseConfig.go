package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type IClickhouseConfig interface {
	HTTPConnStr() string
	NativeConnStr() string
	FlushInterval() int
	FlushLimit() int
}

type ClickhouseConfig struct {
	Viper *viper.Viper
}

func getHost(c *ClickhouseConfig) string {
	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("host")
	c.Viper.SetDefault("host", "0.0.0.0")
	host := c.Viper.GetString("host")
	return host
}

func getDBName(c *ClickhouseConfig) string {
	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("dbname")
	c.Viper.SetDefault("dbname", "webalytic")
	name := c.Viper.GetString("dbname")
	return name
}

func (c ClickhouseConfig) HTTPConnStr() string {
	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("http_port")
	c.Viper.SetDefault("http_port", 8123)
	port := c.Viper.GetInt("http_port")

	host := getHost(&c)
	return fmt.Sprintf("http://%s:%d", host, port)
}

func (c ClickhouseConfig) NativeConnStr() string {
	host := getHost(&c)
	dbname := getDBName(&c)

	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("service_port")
	c.Viper.SetDefault("service_port", 9000)
	port := c.Viper.GetInt("service_port")
	return fmt.Sprintf("http://%s:%d?database=%s", host, port, dbname)
}

func (c ClickhouseConfig) FlushInterval() int {
	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("flush_interval")
	c.Viper.SetDefault("flush_interval", 100)
	v := c.Viper.GetInt("flush_interval")
	return v
}

func (c ClickhouseConfig) FlushLimit() int {
	c.Viper.SetEnvPrefix("CLICKHOUSE")
	c.Viper.BindEnv("flush_limit")
	c.Viper.SetDefault("flush_limit", 100)
	v := c.Viper.GetInt("flush_limit")
	return v
}
