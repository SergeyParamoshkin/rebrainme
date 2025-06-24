package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	DatabaseURL = "postgres://usr:pwd@127.0.0.1:5432/example?sslmode=disable"
)

type slogWrapper struct {
	logger *slog.Logger
}

// Error logs a message at error priority
func (w *slogWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *slogWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Info(msg, args...)
}

func newTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-service"),
			semconv.ServiceVersionKey.String("1.0.0"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	ctx := context.Background()

	// Инициализация slog логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	tp, err := newTracer()
	if err != nil {
		logger.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down tracer provider", "error", err)
		}
	}()

	tracer := otel.Tracer("example-tracer")

	a := app{}

	if err := a.New(ctx, logger, tracer); err != nil {
		logger.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}

	if err := a.Serve(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
