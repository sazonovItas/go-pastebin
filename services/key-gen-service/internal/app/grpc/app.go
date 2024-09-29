package grpcapp

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sazonovItas/go-pastebin/pkg/logger"
	grpchandler "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/handler/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	APITokenMeta = "x-api-key-gen-token"
)

type grpcApp struct {
	log *zap.Logger
	cfg GRPCServerConfig

	grpcServer *grpc.Server
}

func NewApp(
	log *zap.Logger,
	cfg GRPCServerConfig,
	keyGenSvc grpchandler.KeyGenSvc,
) *grpcApp {
	loggingOpts := logging.UnaryServerInterceptor(
		logger.GRPCInterceptor(log),
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			loggingOpts,
			recovery.UnaryServerInterceptor(),
			loggerToContext(log),
		),
	}

	srv := grpc.NewServer(opts...)
	grpchandler.Register(srv, grpchandler.NewKeyHandler(keyGenSvc))

	return &grpcApp{
		log: log,
		cfg: cfg,

		grpcServer: srv,
	}
}

func (a *grpcApp) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *grpcApp) Run(ctx context.Context) error {
	const op = "grpcapp.Run"

	a.log.Info("Starting grpc server", zap.String("port", a.cfg.Port))

	listener, err := net.Listen("tcp", a.cfg.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *grpcApp) Shutdown() {
	a.log.Info("Stopping grpc server")

	a.grpcServer.GracefulStop()
}

func loggerToContext(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return handler(logger.ToContext(ctx, l), req)
	}
}
