package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	httpErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)
)

var (
	memoryUsageGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	cpuUsageGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage in percent",
		},
	)
)

var requestLatencyHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_latency_seconds",
		Help:    "Request latency in seconds",
		Buckets: []float64{0.00001, 0.2, 0.5, 1.0, 2.0, 5.0},
	},
	[]string{"method", "path"},
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		time.Sleep(time.Duration(time.Duration(rand.Int64N(10)).Seconds()))
		latency := time.Since(start).Seconds()

		log.Println(latency)

		httpRequestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()

		requestLatencyHistogram.WithLabelValues(r.Method, r.URL.Path).Observe(latency)
	})
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	prometheus.MustRegister(httpRequestCounter, httpErrorCounter, memoryUsageGauge, cpuUsageGauge, requestLatencyHistogram)

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.Handle("/", metricsMiddleware(http.HandlerFunc(exampleHandler)))

	go collectSystemMetrics()

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func collectSystemMetrics() {
	for {

		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)
		memoryUsageGauge.Set(float64(memStats.Alloc))
		cpuUsageGauge.Set(float64(rand.Int64()))

		time.Sleep(time.Second)
	}
}
