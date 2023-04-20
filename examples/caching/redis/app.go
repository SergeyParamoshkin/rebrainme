package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type app struct {
	logger *zap.Logger

	pool       *pgxpool.Pool
	repository Repository
}

func (a *app) parseUserID(r *http.Request) (*uuid.UUID, error) {
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

		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("userID parsed: %s", userID))

	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("usersHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	ctx := r.Context()

	users, err := a.repository.GetUsers(ctx)
	if err != nil {
		msg := fmt.Sprintf(`failed to get users: %s`, err)

		a.logger.Error(msg)

		writeResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeJsonResponse(w, http.StatusOK, users)
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("userHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	ctx := r.Context()

	userID, err := a.parseUserID(r)
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
		}

		writeResponse(w, status, fmt.Sprintf(`failed to get user with id %s: %s`, userID, err))
		return
	}

	writeJsonResponse(w, http.StatusOK, user)
}

func (a *app) userArticlesHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("userArticlesHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	userID, err := a.parseUserID(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user's id: %s`, err))
		return
	}

	articles, err := a.repository.GetUserArticles(r.Context(), *userID)
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

func (a *app) Init(ctx context.Context, logger *zap.Logger) error {
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
	a.pool = pool
	a.repository = NewCachedRepository(NewRepository(a.pool))

	return a.repository.InitSchema(ctx)
}

func (a *app) Serve() error {
	r := chi.NewRouter()

	r.Get("/users", http.HandlerFunc(a.usersHandler))
	r.Get("/users/{id}", http.HandlerFunc(a.userHandler))
	r.Get("/users/{id}/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/panic", http.HandlerFunc(a.panicHandler))

	// profiling
	r.Mount("/debug", Profiler())

	return http.ListenAndServe("0.0.0.0:9000", r)
}

func Profiler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/cpu", pprof.Handler("cpu"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	return r
}
