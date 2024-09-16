package app

import (
	"context"
	"errors"
	"sync/atomic"

	gatewayapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/gateway"
	grpcapp "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app/grpc"
	"github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/lib/generator"
	keygensvc "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/service/keygen"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type (
	Service interface {
		Run(ctx context.Context) error
		Shutdown()
	}
)

type app struct {
	cfg Config

	ctx    context.Context
	cancel context.CancelFunc

	services    []Service
	cleanUps    []func()
	cleanUpDone atomic.Int32
}

func NewApp(log *zap.Logger, cfg Config) *app {
	authFn := func(token string) error {
		if token != cfg.Core.APIToken {
			return errors.New("invalid api token")
		}

		return nil
	}

	keyGenService := keygensvc.New(generator.NewGenerator(cfg.Core.KeyLength), cfg.Core.KeyBuffer)
	grpcServer := grpcapp.NewApp(log.Named("grpc_server"), cfg.GRPCServer, keyGenService, authFn)
	gatewayServer := gatewayapp.NewApp(
		log.Named("gateway_server"),
		cfg.GatewayServer,
		keyGenService,
	)

	return &app{
		cfg: cfg,

		services: []Service{grpcServer, gatewayServer},
		cleanUps: []func(){keyGenService.Stop},
	}
}

func (a *app) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *app) Run(parentCtx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(parentCtx)

	wg, ctx := errgroup.WithContext(a.ctx)
	for _, svc := range a.services {
		wg.Go(func() error {
			return svc.Run(ctx)
		})
	}
	<-ctx.Done()

	a.shutdownServices()

	return wg.Wait()
}

func (a *app) Shutdown() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *app) CleanUp() {
	if !a.cleanUpDone.CompareAndSwap(0, 1) {
		return
	}

	for _, cleanUp := range a.cleanUps {
		cleanUp()
	}
}

func (a *app) shutdownServices() {
	for _, svc := range a.services {
		svc.Shutdown()
	}
}
