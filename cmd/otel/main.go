package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() 
	exporter := prometheus.New()
	if err != nil {
		log.Fatalf("failed to initialize prometheus exporter: %v", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	meter := otel.Meter("example-meter")

	counter, err := meter.Int64Counter("example_counter", otelmetric.WithDescription("A simple counter"))
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	upDownCounter, err := meter.Int64UpDownCounter("example_updown_counter", otelmetric.WithDescription("A simple up-down counter"))
	if err != nil {
		log.Fatalf("failed to create up-down counter: %v", err)
	}

	histogram, err := meter.Float64Histogram("example_histogram", otelmetric.WithDescription("A simple histogram"))
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	gauge, err := meter.Float64ObservableGauge("example_gauge", otelmetric.WithDescription("A simple gauge"))
	if err != nil {
		log.Fatalf("failed to create gauge: %v", err)
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, o otelmetric.Observer) error {
		o.ObserveFloat64(gauge, rand.Float64()*100)
		return nil
	}, gauge)
	if err != nil {
		log.Fatalf("failed to register callback: %v", err)
	}

	go func() {
		http.Handle("/metrics", exporter)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	ctx := context.Background()
	for {

		counter.Add(ctx, 1)
		if rand.Intn(2) == 0 {
			upDownCounter.Add(ctx, 1)
		} else {
			upDownCounter.Add(ctx, -1)
		}

		histogram.Record(ctx, rand.Float64()*100)

		time.Sleep(time.Second)
	}
}
