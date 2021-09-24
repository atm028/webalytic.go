package Consul

import (
	"fmt"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/hashicorp/consul/api"
	CommonCfg "github.com/webalytic.go/common/config"
	"go.uber.org/fx"

	"net"
	"strings"
)

type IConsul interface {
	ServiceRegister(
		name string,
		port int,
		interval string,
		path string,
		method string,
	) bool
}

type Consul struct {
	client *api.Client
	agent  *api.Agent
	logger *bunyan.Logger
}

func (c *Consul) ServiceRegister(
	name string,
	port int,
	interval string,
	path string,
	method string) bool {
	if c == nil {
		return false
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		c.logger.Debug("cant get net addresses")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddr := strings.Split(ipnet.String(), "/")[0]
				checkUri := fmt.Sprintf("http://%s:%d%s", ipAddr, port, path)
				c.logger.Debug(fmt.Sprintf("Register Consul check URI %s", checkUri))
				svc := &api.AgentServiceRegistration{
					Name:    "collector",
					ID:      name,
					Port:    port,
					Address: ipAddr,
					//Name: name,
					Check: &api.AgentServiceCheck{
						Interval: interval,
						HTTP:     checkUri,
						Method:   method,
						Status:   "passing",
					},
				}
				if err := c.agent.ServiceRegister(svc); err != nil {
					c.logger.Fatal(fmt.Sprintf("Consul service registration error: %s", err))
					return false
				}
			}
		}
	}
	return true
}

func Container() fx.Option {
	return fx.Options(fx.Provide(func(
		logger bunyan.Logger,
		cfg *CommonCfg.ConsulConfig,
	) *Consul {
		if !cfg.Enabled() {
			logger.Debug("Consul is disabled")
			return nil
		}

		consulConfig := cfg.Config()
		consulClient, _ := api.NewClient(consulConfig)
		return &Consul{
			client: consulClient,
			agent:  consulClient.Agent(),
			logger: &logger,
		}
	}))
}
