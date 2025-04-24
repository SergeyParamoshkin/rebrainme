package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	trace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type app struct {
	logger *zap.Logger
	tracer trace.Tracer

	pool       *pgxpool.Pool
	repository *Repository
}

func (a *app) readyHandler(w http.ResponseWriter, r *http.Request) {
	err := a.repository.pool.Ping(r.Context())
	if err != nil {
		a.logger.Error(err.Error())
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// cache
	// .... 

	// other deps
	// .... 
	//

	writeResponse(w, http.StatusOK, "ready")
}

func (a *app) parseUserID(ctx context.Context, r *http.Request) (*uuid.UUID, error) {
	ctx, span := a.tracer.Start(ctx, "parseUserID")
	defer span.End()

	strUserID := chi.URLParam(r, "id")

	if strUserID == "" {
		return nil, nil
	}

	userID, err := uuid.Parse(strUserID)
	if err != nil {
		a.logger.Debug(
			fmt.Sprintf("failed to parse userID (uuid) from: '%s'", strUserID),
			zap.Field{Key: "error", String: err.Error(), Type: zapcore.StringType},
		)

		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("MyErrorType"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))

		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("userID parsed: %s", userID))

	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.usersHandler")
	defer span.End()

	a.logger.Info("usersHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	users, err := a.repository.GetUsers(ctx)
	if err != nil {
		msg := fmt.Sprintf(`failed to get users: %s`, err)

		a.logger.Error(msg)

		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("MyErrorType"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))

		writeResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeJsonResponse(w, http.StatusOK, users)
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.userHandler")
	defer span.End()

	a.logger.Info("userHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	userID, err := a.parseUserID(ctx, r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user's id: %s`, err))
		return
	}

	user, err := a.repository.GetUser(ctx, *userID)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, &NotFoundError{}):
			hub := sentry.GetHubFromContext(r.Context())
			hub.CaptureException(err)
			status = http.StatusNotFound
			span.SetStatus(codes.Error, "Not Found")
			span.AddEvent("not_found", trace.WithAttributes(
				semconv.ExceptionTypeKey.String("*NotFoundError"),
				semconv.ExceptionMessageKey.String(err.Error()),
			))

		default:
			span.AddEvent("error", trace.WithAttributes(
				semconv.ExceptionTypeKey.String("MyErrorType"),
				semconv.ExceptionMessageKey.String(err.Error()),
			))
		}

		writeResponse(w, status, fmt.Sprintf(`failed to get user with id %s`, userID))
		return
	}

	writeJsonResponse(w, http.StatusOK, user)
}

func (a *app) userArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.userArticlesHandler")
	defer span.End()

	a.logger.Info("userArticlesHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	userID, err := a.parseUserID(ctx, r)
	if err != nil {
		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("MyErrorType"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))

		writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user's id: %s`, err))
		return
	}

	articles, err := a.repository.GetUserArticles(ctx, *userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf(`failed to get user's (id: %s) articles: %s`, userID, err))
		return
	}

	writeJsonResponse(w, http.StatusOK, articles)
}

func (a *app) panicHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = recover()

		writeResponse(w, http.StatusInternalServerError, "panic logged, see server log")
	}()

	a.logger.Panic("panic!!!")
}

func (a *app) New(ctx context.Context, logger *zap.Logger, tracer trace.Tracer) error {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse conn string (%s): %w", DatabaseURL, err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	config.ConnConfig.Logger = zapadapter.NewLogger(logger) // логгер запросов в БД

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	a.logger = logger
	a.tracer = tracer
	a.pool = pool
	a.repository = NewRepository(a.pool, a.tracer)

	return a.repository.InitSchema(ctx)
}

func (a *app) Serve() error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   "http://8a8c6c3d8d7aebf6ef0fd39ca2e9fdfa@localhost:9000/3",
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(time.Second)

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "chi-http-server")
	})
	r.Use(middleware.Recoverer)
	r.Use(sentryMiddleware.Handle)

	r.Get("/users", http.HandlerFunc(a.usersHandler))
	r.Get("/users/{id}", http.HandlerFunc(a.userHandler))
	r.Get("/users/{id}/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		hub := sentry.GetHubFromContext(r.Context())
		hub.CaptureException(errors.New("test error"))
	})
	r.Get("/panic", http.HandlerFunc(a.panicHandler))

	// Liveness probe (проверка, что сервис "жив")
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Readiness probe (проверка, что сервис готов принимать трафик)
	r.Get("/ready", http.HandlerFunc(a.readyHandler))

	return http.ListenAndServe("0.0.0.0:8000", r)
}
