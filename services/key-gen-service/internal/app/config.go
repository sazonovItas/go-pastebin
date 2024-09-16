package app

import (
	gatewayapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/gateway"
	grpcapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/grpc"
)

type Config struct {
	Log           LogConfig                      `yaml:"log"     mapstructure:"log"`
	Core          CoreConfig                     `yaml:"core"    mapstructure:"core"`
	GRPCServer    grpcapp.GRPCServerConfig       `yaml:"grpc"    mapstructure:"grpc"`
	GatewayServer gatewayapp.GatewayServerConfig `yaml:"gateway" mapstructure:"gateway"`
}

type LogConfig struct {
	Level  string `yaml:"level"  mapstructure:"level"`
	Format string `yaml:"format" mapstructure:"format"`
}

type CoreConfig struct {
	KeyBuffer int    `yaml:"key_buffer" mapstructure:"key_buffer"`
	KeyLength int    `yaml:"key_length" mapstructure:"key_length"`
	APIToken  string `yaml:"api_token"  mapstructure:"api_token"`
}
