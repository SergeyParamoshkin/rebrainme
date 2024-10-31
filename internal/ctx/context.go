package fwctx

import (
	"context"

	"go.uber.org/zap"
)

type key uint8

const (
	KeyLogger key = iota
	KeyErrorInterceptor
	KeyToken
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, KeyLogger, logger)
}

func LoggerFromCtx(ctx context.Context, fallbackLogger *zap.Logger) *zap.Logger {
	logger, _ := ctx.Value(KeyLogger).(*zap.Logger)

	if logger == nil {
		return fallbackLogger
	}

	return logger
}

func Empty(logger *zap.Logger) context.Context {
	return WithLogger(context.Background(), logger)
}
