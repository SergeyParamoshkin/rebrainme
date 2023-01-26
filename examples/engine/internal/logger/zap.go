package logger

import (
	"context"

	"go.uber.org/zap"
)

type zapLoggerKeyType string

const zapLoggerKey zapLoggerKeyType = "zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: logger,
	}
}

func (l *ZapLogger) Debug(err error, tags ...Tag) {
	l.logger.Debug(err.Error(), makeFields(tags)...)
}

func (l *ZapLogger) Info(err error, tags ...Tag) {
	l.logger.Info(err.Error(), makeFields(tags)...)
}

func (l *ZapLogger) Error(err error, tags ...Tag) {
	l.logger.Error(err.Error(), makeFields(tags)...)
}

func (l ZapLogger) With(tags ...Tag) Logger {
	l.logger = l.logger.With(makeFields(tags)...)

	return &l
}

func (l ZapLogger) NewContext(ctx context.Context, tags ...Tag) context.Context {
	return context.WithValue(ctx, zapLoggerKey, l.WithContext(ctx).With(tags...))
}

func (l ZapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return &l
	}

	if ctxLogger, ok := ctx.Value(zapLoggerKey).(ZapLogger); ok {
		return &ctxLogger
	}

	return &l
}

func makeFields(tags []Tag) []zap.Field {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}

	return fields
}
