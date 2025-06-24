package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	trace "go.opentelemetry.io/otel/trace"
)

type retryCountKey int64

// SlogAdapter реализует интерфейс pgx.Logger
type SlogAdapter struct {
	logger *slog.Logger
}

func (a *SlogAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var slogLevel slog.Level
	switch level {
	case pgx.LogLevelTrace:
		slogLevel = slog.LevelDebug - 2
	case pgx.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case pgx.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case pgx.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case pgx.LogLevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	a.logger.Log(ctx, slogLevel, msg, "pgx", data)
}

type app struct {
	logger *slog.Logger
	tracer trace.Tracer
	pool   *pgxpool.Pool
	repo   *Repository
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
		a.logger.DebugContext(ctx, "failed to parse userID",
			"input", strUserID,
			"error", err,
		)

		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("ParseError"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))

		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	a.logger.DebugContext(ctx, "userID parsed", "userID", userID)
	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.usersHandler")
	defer span.End()

	a.logger.InfoContext(ctx, "usersHandler called", "method", r.Method)

	a.logger.WithGroup("request").Debug("API.usersHandler",
		"method", r.Method,
		"url", r.URL)

	const retryCountKey = "retryCount"
	ctx = context.WithValue(ctx, retryCountKey, 10)

	users, err := a.repo.GetUsers(ctx)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to get users", "error", err)
		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("DatabaseError"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))
		writeResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJsonResponse(w, http.StatusOK, users)
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.userHandler")
	defer span.End()

	a.logger.InfoContext(ctx, "userHandler called", "method", r.Method)

	userID, err := a.parseUserID(ctx, r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.repo.GetUser(ctx, *userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		} else {
			span.AddEvent("error", trace.WithAttributes(
				semconv.ExceptionTypeKey.String("DatabaseError"),
				semconv.ExceptionMessageKey.String(err.Error()),
			))
		}
		writeResponse(w, status, err.Error())
		return
	}

	writeJsonResponse(w, http.StatusOK, user)
}

func (a *app) userArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "API.userArticlesHandler")
	defer span.End()

	a.logger.InfoContext(ctx, "userArticlesHandler called", "method", r.Method)

	userID, err := a.parseUserID(ctx, r)
	if err != nil {
		span.AddEvent("error", trace.WithAttributes(
			semconv.ExceptionTypeKey.String("ParseError"),
			semconv.ExceptionMessageKey.String(err.Error()),
		))
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	articles, err := a.repo.GetUserArticles(ctx, *userID)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to get articles",
			"userID", userID,
			"error", err,
		)
		writeResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJsonResponse(w, http.StatusOK, articles)
}

func (a *app) panicHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			a.logger.ErrorContext(r.Context(), "panic occurred",
				"error", err,
			)
			writeResponse(w, http.StatusOK, "panic handled")
		}
	}()

	a.logger.ErrorContext(r.Context(), "panic triggered")
	panic("manual panic")
}

func (a *app) New(ctx context.Context, logger *slog.Logger, tracer trace.Tracer) error {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	config.ConnConfig.Logger = &SlogAdapter{logger: logger}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	a.logger = logger
	a.tracer = tracer
	a.pool = pool
	a.repo = NewRepository(pool, tracer, logger)

	return a.repo.InitSchema(ctx)
}

func (a *app) Serve() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Route("/debug", func(r chi.Router) {
		r.Mount("/", middleware.Profiler())
	})

	r.Get("/users", a.usersHandler)
	r.Get("/users/{id}", a.userHandler)
	r.Get("/users/{id}/articles", a.userArticlesHandler)
	r.Get("/panic", a.panicHandler)

	// r.Route("/debug/pprof", func(p chi.Router) {
	// 	p.Get("/", pprof.Index)
	// 	p.Get("/cmdline", pprof.Cmdline)
	// 	p.Get("/profile", pprof.Profile)
	// 	p.Post("/symbol", pprof.Symbol)
	// 	p.Get("/symbol", pprof.Symbol)
	// 	p.Get("/trace", pprof.Trace)
	// 	p.Get("/{profile}", pprof.Index)
	// })

	return http.ListenAndServe(":9000", r)
}
