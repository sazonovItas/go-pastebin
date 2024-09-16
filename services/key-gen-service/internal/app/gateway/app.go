package gatewayapp

import (
	"context"

	grpchandler "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/handler/grpc"
	"go.uber.org/zap"
)

type gatewayApp struct {
	log *zap.Logger
	cfg GatewayServerConfig
}

func NewApp(log *zap.Logger, cfg GatewayServerConfig, keyGenSvc grpchandler.KeyGenSvc) *gatewayApp {
	return &gatewayApp{
		log: log,
		cfg: cfg,
	}
}

func (a *gatewayApp) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *gatewayApp) Run(ctx context.Context) error {
	const op = "grpcapp.Run"

	a.log.Info("Starting gateway server")

	return nil
}

func (a *gatewayApp) Shutdown() {
	a.log.Info("Stopping gateway server")
}
