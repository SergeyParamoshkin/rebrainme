package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergeyParamoshkin/rebrainme/internal/api"
	"github.com/SergeyParamoshkin/rebrainme/internal/logger"
	customMiddleware "github.com/SergeyParamoshkin/rebrainme/internal/middleware"
	"github.com/SergeyParamoshkin/rebrainme/internal/monitoring"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		debug = flag.Bool("debug", false, "Enable debug mode")
		port  = flag.String("port", "8080", "Server port")
	)
	flag.Parse()

	// Инициализация логгера
	logger.Init(*debug)
	logger.Info("Starting application",
		"port", *port,
		"debug", *debug,
	)

	// Инициализация метрик
	monitoring.InitMetrics()

	// Создаём chi роутер
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(httplog.NewLogger("api", httplog.Options{
		JSON:             !*debug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
	})))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Наш middleware
	rateLimiter := customMiddleware.NewRateLimiter(10, 100)
	router.Use(customMiddleware.RateLimitMiddleware(rateLimiter))
	router.Use(customMiddleware.MetricsMiddleware)

	// Регистрация маршрутов
	api.RegisterHandlers(router)

	// Prometheus метрики
	router.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server stopped gracefully")
}
