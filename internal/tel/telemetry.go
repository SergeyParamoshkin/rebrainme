package tel

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"

	promCollectors "github.com/prometheus/client_golang/prometheus/collectors"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Telemetry struct {
	logger       *zap.Logger
	registry     *prometheus.Registry
	tracer       opentracing.Tracer
	tracerCloser io.Closer
}

func New(logger *zap.Logger, config *Config) (*Telemetry, error) {
	registry := prometheus.NewRegistry()

	{ // default collectors
		pc := promCollectors.NewProcessCollector(promCollectors.ProcessCollectorOpts{
			Namespace: config.ServiceName,
		})

		if err := registry.Register(pc); err != nil {
			return nil, fmt.Errorf("failed to register process collector: %w", err)
		}

		if err := registry.Register(promCollectors.NewGoCollector()); err != nil {
			return nil, fmt.Errorf("failed to register go collector: %w", err)
		}
	}

	cfg := jaegercfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeRemote,
			Param: config.SamplerFraction,
		},
		// TODO: Add reporter!
	}

	tracer, tracerCloser, err := cfg.NewTracer()
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger tracer: %w", err)
	}

	return &Telemetry{
		logger:       logger,
		registry:     registry,
		tracer:       tracer,
		tracerCloser: tracerCloser,
	}, nil
}

func (t *Telemetry) Close() error {
	return t.tracerCloser.Close()
}

func (t *Telemetry) Registry() *prometheus.Registry {
	return t.registry
}

func (t *Telemetry) Register(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := t.registry.Register(collector); err != nil {
			return fmt.Errorf("failed to register collector: %w", err)
		}
	}

	return nil
}

func (t *Telemetry) StartSpan( //nolint:ireturn // opentracing
	ctx context.Context,
	spanName string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, t.tracer, spanName, opts...,
	)

	return ctx, span
}
