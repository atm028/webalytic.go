package config

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type IConsulConfig interface {
	Enabled() bool
	Config() string
}

type ConsulConfig struct {
	Viper *viper.Viper
}

func (c ConsulConfig) address() string {
	c.Viper.SetEnvPrefix("CONSUL")
	c.Viper.BindEnv("address")
	c.Viper.SetDefault("addr", "0.0.0.0:8500")
	addr := c.Viper.GetString("address")
	return addr
}

func (c ConsulConfig) Config() *api.Config {
	cfg := api.DefaultConfig()
	cfg.Address = c.address()
	return cfg
}

func (c ConsulConfig) Enabled() bool {
	c.Viper.SetEnvPrefix("CONSUL")
	c.Viper.BindEnv("enabled")
	c.Viper.SetDefault("enabled", false)
	v := c.Viper.GetBool("enabled")
	return v
}
