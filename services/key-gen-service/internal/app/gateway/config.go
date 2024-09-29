package gatewayapp

type GatewayServerConfig struct {
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
	Port    string `yaml:"port"    mapstructure:"port"`
}
