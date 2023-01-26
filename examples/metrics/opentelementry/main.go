package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	KeyMethod = attribute.Key("method")
	KeyStatus = attribute.Key("status")
)

type app struct {
	attrs              []attribute.KeyValue
	latencyMsRecorder  syncfloat64.Histogram
	lineLengthRecorder syncint64.Histogram
	lineCounter        syncint64.Counter
	lastLineLength     syncint64.UpDownCounter
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()
	commonLabels := []attribute.KeyValue{KeyMethod.String(r.Method), KeyStatus.String("OK")}
	a.attrs = append(a.attrs, commonLabels...)

	line := r.URL.Query().Get("line")
	lineLength := int64(len(line))

	defer func(ctx context.Context) {
		a.latencyMsRecorder.Record(ctx, sinceInMilliseconds(startTime), a.attrs...)
		a.lineLengthRecorder.Record(ctx, lineLength, a.attrs...)
		a.lineCounter.Add(ctx, 1, a.attrs...)
		a.lastLineLength.Add(ctx, lineLength, a.attrs...)
	}(ctx)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // имитация работы
	writeResponse(w, http.StatusOK, strings.ToUpper(line))
}

func (a *app) initMeters(provider *metric.MeterProvider) error {
	var err error

	a.attrs = []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}

	meter := provider.Meter("rebrainmemetrics")

	a.latencyMsRecorder, err = meter.SyncFloat64().
		Histogram("repl/latency", instrument.WithDescription("The distribution of the latencies"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	a.lineLengthRecorder, err = meter.SyncInt64().
		Histogram("repl/line_lengths", instrument.WithDescription("Groups the lengths of keys in buckets"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	a.lineCounter, err = meter.SyncInt64().
		UpDownCounter("repl/line_count", instrument.WithDescription("Count of lines"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	a.lastLineLength, err = meter.SyncInt64().
		UpDownCounter("repl/last_line_length", instrument.WithDescription("Last line length"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return err
}

func (a *app) Init(ctx context.Context) error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	if err := a.initMeters(provider); err != nil {
		return err
	}

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go a.Serve()

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()

	return nil
}

func (a *app) Serve() error {
	log.Printf("serving metrics at http://localhost:9000/metrics")

	mux := http.NewServeMux()
	mux.Handle("/process", http.HandlerFunc(a.processHandler)) // /process?line=текст+тут
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler:           mux,
		Addr:              "0.0.0.0:9000",
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return err
}

func main() {
	ctx := context.Background()
	a := app{}

	if err := a.Init(ctx); err != nil {
		log.Fatal(err)
	}
}

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}

func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}
