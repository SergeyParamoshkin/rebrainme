package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/internal/logger"
	"github.com/SergeyParamoshkin/rebrainme/internal/monitoring"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		endpoint := r.URL.Path
		method := r.Method

		monitoring.ActiveRequests.WithLabelValues(method, endpoint).Inc()
		defer monitoring.ActiveRequests.WithLabelValues(method, endpoint).Dec()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)

		monitoring.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		monitoring.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)

		// Добавляем логирование
		reqID := middleware.GetReqID(r.Context())
		logger.LogRequest(method, endpoint, wrapped.statusCode, time.Since(start),
			"request_id", reqID,
			"remote_addr", r.RemoteAddr,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter doesn't support hijacking")
	}
	return hijacker.Hijack()
}