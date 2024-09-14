package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func GRPCInterceptor(l *zap.Logger) logging.Logger {
	log := l.Sugar()
	return logging.LoggerFunc(
		func(ctx context.Context, level logging.Level, msg string, fields ...any) {
			largs := append([]any{"msg", msg}, fields)
			switch level {
			case logging.LevelDebug:
				log.Debug(largs...)
			case logging.LevelInfo:
				log.Info(largs...)
			case logging.LevelWarn:
				log.Warn(largs...)
			case logging.LevelError:
				log.Error(largs...)
			default:
				log.Info(largs...)
			}
		},
	)
}
