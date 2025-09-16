package middleware

import (
	"net/http"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/internal/monitoring"
	"github.com/SergeyParamoshkin/rebrainme/internal/ratelimit"
)

type RateLimiter struct {
	tokenBucket   *ratelimit.TokenBucket
	slidingWindow *ratelimit.SlidingWindow
}

func NewRateLimiter(tokensPerSecond float64, windowLimit int) *RateLimiter {
	return &RateLimiter{
		tokenBucket:   ratelimit.NewTokenBucket(tokensPerSecond, tokensPerSecond*10),
		slidingWindow: ratelimit.NewSlidingWindow(time.Minute, windowLimit),
	}
}

func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.tokenBucket.Allow() {
				monitoring.RateLimitedRequests.WithLabelValues("token_bucket").Inc()
				http.Error(w, "Rate limit exceeded (token bucket)", http.StatusTooManyRequests)
				return
			}

			if !limiter.slidingWindow.Allow() {
				monitoring.RateLimitedRequests.WithLabelValues("sliding_window").Inc()
				http.Error(w, "Rate limit exceeded (sliding window)", http.StatusTooManyRequests)
				return
			}

			monitoring.TokenBucketAvailable.WithLabelValues("main").Set(float64(limiter.tokenBucket.Available()))
			monitoring.SlidingWindowRequests.WithLabelValues("main").Set(float64(limiter.slidingWindow.Count()))

			next.ServeHTTP(w, r)
		})
	}
}
