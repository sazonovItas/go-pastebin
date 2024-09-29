package logger

import (
	"context"

	"go.uber.org/zap"
)

func Sync() error {
	return getLogger().Sync()
}

func Named(name string) *zap.Logger {
	return getLogger().Named(name)
}

func With(fields ...zap.Field) *zap.Logger {
	return getLogger().With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getLogger().Panic(msg, fields...)
}

func Sugar() *zap.SugaredLogger {
	return getLogger().Sugar()
}

func NamedContext(ctx context.Context, name string) *zap.Logger {
	return FromContext(ctx).Named(name)
}

func WithContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	return FromContext(ctx).With(fields...)
}

func DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Debug(msg, fields...)
}

func InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, fields...)
}

func WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, fields...)
}

func ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, fields...)
}

func FatalContext(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Fatal(msg, fields...)
}

func PanicContetxt(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Panic(msg, fields...)
}

func SugarContext(ctx context.Context) *zap.SugaredLogger {
	return FromContext(ctx).Sugar()
}

func getLogger() *zap.Logger {
	return CreateLogger()
}
