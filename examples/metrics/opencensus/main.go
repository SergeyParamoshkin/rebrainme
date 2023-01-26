package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	KeyMethod, _ = tag.NewKey("method")
	KeyStatus, _ = tag.NewKey("status")
)

type app struct {
	pe *prometheus.Exporter

	MLatencyMs   *stats.Float64Measure
	MLineLengths *stats.Int64Measure

	latencyView,
	lineCountView,
	lineLengthView,
	lastLineLengthView *view.View
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	ctx, err := tag.New(
		r.Context(),
		tag.Insert(KeyMethod, r.Method),
		tag.Insert(KeyStatus, "OK"))

	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("tag error: %s", err))
		return
	}

	line := r.URL.Query().Get("line")

	defer func() {
		stats.Record(
			ctx,
			a.MLatencyMs.M(sinceInMilliseconds(startTime)),
			a.MLineLengths.M(int64(len(line))))
	}()

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // имитация работы

	writeResponse(w, http.StatusOK, strings.ToUpper(line))
}

func (a *app) Init() error {
	// время обработки в мс
	a.MLatencyMs = stats.Float64("repl/latency", "The latency in milliseconds per REPL loop", "ms")

	// распределение длин строк
	a.MLineLengths = stats.Int64("repl/line_lengths", "The distribution of line lengths", "By")

	// prometheus type: histogram
	a.latencyView = &view.View{
		Name:        "demo/latency",
		Measure:     a.MLatencyMs,
		Description: "The distribution of the latencies",
		// границы гистограммы
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		TagKeys:     []tag.Key{KeyMethod}}

	// prometheus type: counter
	a.lineCountView = &view.View{
		Name:        "demo/lines_in",
		Measure:     a.MLineLengths,
		Description: "The number of lines from standard input",
		Aggregation: view.Count(),
	}

	// prometheus type: histogram
	a.lineLengthView = &view.View{
		Name:        "demo/line_lengths",
		Description: "Groups the lengths of keys in buckets",
		Measure:     a.MLineLengths,
		// длины: [>=0B, >=5B, >=10B, >=15B, >=20B, >=40B, >=60B, >=80, >=100B, >=200B, >=400, >=600, >=800, >=1000]
		Aggregation: view.Distribution(0, 5, 10, 15, 20, 40, 60, 80, 100, 200, 400, 600, 800, 1000),
	}

	// prometheus type: gauge
	a.lastLineLengthView = &view.View{
		Name:        "demo/last_line_length",
		Measure:     a.MLineLengths,
		Description: "The length of last line",
		Aggregation: view.LastValue(),
	}

	err := view.Register(a.latencyView, a.lineCountView, a.lineLengthView, a.lastLineLengthView)
	if err != nil {
		return err
	}

	a.pe, err = prometheus.NewExporter(prometheus.Options{
		Namespace: "ocmetricsexample",
	})

	return err
}

func (a *app) Serve() error {
	mux := http.NewServeMux()
	mux.Handle("/process", http.HandlerFunc(a.processHandler)) // /process?line=текст+тут
	mux.Handle("/metrics", a.pe)

	return http.ListenAndServe("0.0.0.0:9000", mux)
}

func main() {
	a := app{}

	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	if err := a.Serve(); err != nil {
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
