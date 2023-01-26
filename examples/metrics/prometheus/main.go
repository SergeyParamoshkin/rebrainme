package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Namespace = "ocmetricsexample"

	LabelMethod = "method"
	LabelStatus = "status"
)

type app struct {
	latencyHistogram,
	lineLengthHistogram *prometheus.HistogramVec

	lineCounter prometheus.Counter

	lastLineLengthGauge prometheus.Gauge
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	line := r.URL.Query().Get("line")

	defer func() {
		a.latencyHistogram.With(prometheus.Labels{LabelMethod: r.Method}).Observe(sinceInMilliseconds(startTime))

		a.lineLengthHistogram.With(prometheus.Labels{LabelStatus: "OK"}).Observe(float64(len(line)))

		a.lineCounter.Inc()

		a.lastLineLengthGauge.Set(float64(len(line)))
	}()

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // имитация работы

	writeResponse(w, http.StatusOK, strings.ToUpper(line))
}

func (a *app) Init() error {
	// prometheus type: histogram
	a.latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "latency",
		Help:      "The distribution of the latencies",
		Buckets:   []float64{0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000},
	}, []string{LabelMethod})

	// prometheus type: histogram
	a.lineLengthHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "line_lengths",
		Help:      "Groups the lengths of keys in buckets",
		// длины: [>=0B, >=5B, >=10B, >=15B, >=20B, >=40B, >=60B, >=80, >=100B, >=200B, >=400, >=600, >=800, >=1000]
		Buckets: []float64{0, 5, 10, 15, 20, 40, 60, 80, 100, 200, 400, 600, 800, 1000},
	}, []string{LabelStatus})

	// prometheus type: counter
	a.lineCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "lines_in",
		Help:      "The number of lines from standard input",
	})

	// prometheus type: gauge
	a.lastLineLengthGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "last_line_length",
		Help:      "The length of last line",
	})

	prometheus.MustRegister(a.latencyHistogram)
	prometheus.MustRegister(a.lineLengthHistogram)
	prometheus.MustRegister(a.lineCounter)
	prometheus.MustRegister(a.lastLineLengthGauge)

	return nil
}

func (a *app) Serve() error {
	mux := http.NewServeMux()
	mux.Handle("/process", http.HandlerFunc(a.processHandler)) // /process?line=текст+тут
	mux.Handle("/metrics", promhttp.Handler())

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
