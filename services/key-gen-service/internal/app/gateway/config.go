package gatewayapp

import "time"

type GatewayServerConfig struct {
	Address string        `yaml:"address" mapstructure:"address"`
	Timeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
}
