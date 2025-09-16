package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	RateLimitedRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limited_requests_total",
			Help: "Total number of rate limited requests",
		},
		[]string{"limiter_type"},
	)

	RetryAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "retry_attempts_total",
			Help: "Total number of retry attempts",
		},
		[]string{"operation", "result"},
	)

	RetryBudgetExhausted = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "retry_budget_exhausted_total",
			Help: "Total number of times retry budget was exhausted",
		},
		[]string{"operation"},
	)

	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of active HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
		},
		[]string{"service"},
	)

	TokenBucketAvailable = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_bucket_available",
			Help: "Available tokens in token bucket",
		},
		[]string{"bucket_name"},
	)

	SlidingWindowRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sliding_window_current_requests",
			Help: "Current requests in sliding window",
		},
		[]string{"window_name"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}