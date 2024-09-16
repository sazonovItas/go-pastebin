package grpcapp

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/sazonovItas/go-pastebin/pkg/logger"
	grpchandler "github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/handler/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	authFn func(token string) error,
) *grpcApp {
	authFunc := func(ctx context.Context) (context.Context, error) {
		md := metadata.ExtractIncoming(ctx)

		token, ok := md[APITokenMeta]
		if !ok {
			return ctx, status.Errorf(codes.Unauthenticated, "missing api key gen token")
		}

		if len(token) != 1 {
			return ctx, status.Errorf(codes.Unauthenticated, "invalid api key gen token format")
		}

		if err := authFn(token[0]); err != nil {
			return ctx, status.Errorf(codes.Unauthenticated, "%s", err.Error())
		}

		return ctx, nil
	}

	loggingOpts := logging.UnaryServerInterceptor(
		logger.GRPCInterceptor(log),
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			loggingOpts,
			recovery.UnaryServerInterceptor(),
			auth.UnaryServerInterceptor(authFunc),
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

	a.log.Info("Starting grpc server", zap.String("address", a.cfg.Address))

	listener, err := net.Listen("tcp", a.cfg.Address)
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
