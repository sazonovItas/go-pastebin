package grpcapp

import "time"

type GRPCServerConfig struct {
	Address string        `yaml:"address" mapstructure:"address"`
	Timeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
}
