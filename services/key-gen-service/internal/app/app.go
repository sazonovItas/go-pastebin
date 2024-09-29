package app

import (
	"github.com/sazonovItas/go-pastebin/pkg/app"
	gatewayapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/gateway"
	grpcapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/grpc"
	"github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/lib/generator"
	keygensvc "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/service/keygen"
	"go.uber.org/zap"
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

func New(log *zap.Logger, cfg Config) *app.App {
	keyGenService := keygensvc.New(generator.NewGenerator(cfg.Core.KeyLength), cfg.Core.KeyBuffer)

	apps := make([]app.Service, 0)
	grpcServer := grpcapp.NewApp(
		log.Named("grpc_server"),
		cfg.GRPCServer,
		keyGenService,
	)
	apps = append(apps, grpcServer)

	if cfg.GatewayServer.Enabled {
		gatewayServer := gatewayapp.NewApp(
			log.Named("gateway_server"),
			cfg.GatewayServer,
			keyGenService,
		)
		apps = append(apps, gatewayServer)
	}

	return &app.App{
		Cfg: cfg,

		Services: apps,
		CleanUps: []func(){keyGenService.Stop},
	}
}
