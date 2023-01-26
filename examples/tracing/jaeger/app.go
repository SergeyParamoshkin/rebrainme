package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

type app struct {
	logger *zap.Logger
	tracer opentracing.Tracer

	pool       *pgxpool.Pool
	repository *Repository
}

func (a *app) parseUserID(ctx context.Context, r *http.Request) (*uuid.UUID, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, a.tracer, "parseUserID")
	defer span.Finish()

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

		span.LogFields(
			log.Error(err),
		)

		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("userID parsed: %s", userID))

	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "usersHandler")
	defer span.Finish()

	a.logger.Info("usersHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	users, err := a.repository.GetUsers(ctx)
	if err != nil {
		msg := fmt.Sprintf(`failed to get users: %s`, err)

		a.logger.Error(msg)

		span.LogFields(
			log.Error(err),
		)

		writeResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeJsonResponse(w, http.StatusOK, users)
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "userHandler")
	defer span.Finish()

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
		case errors.Is(err, ErrNotFound):
			status = http.StatusNotFound
		default:
			span.LogFields(
				log.Error(err),
			)
		}

		writeResponse(w, status, fmt.Sprintf(`failed to get user with id %s: %s`, userID, err))
		return
	}

	writeJsonResponse(w, http.StatusOK, user)
}

func (a *app) userArticlesHandler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "userArticlesHandler")
	defer span.Finish()

	a.logger.Info("userArticlesHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	userID, err := a.parseUserID(ctx, r)
	if err != nil {
		span.LogFields(
			log.Error(err),
		)

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

		writeResponse(w, http.StatusOK, "panic logged, see server log")
	}()

	a.logger.Panic("panic!!!")
}

func (a *app) Init(ctx context.Context, logger *zap.Logger, tracer opentracing.Tracer) error {
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
	r := chi.NewRouter()

	r.Get("/users", http.HandlerFunc(a.usersHandler))
	r.Get("/users/{id}", http.HandlerFunc(a.userHandler))
	r.Get("/users/{id}/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/panic", http.HandlerFunc(a.panicHandler))

	return http.ListenAndServe("0.0.0.0:9000", nethttp.Middleware(a.tracer, r))
}
