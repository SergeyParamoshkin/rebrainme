package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

type app struct {
	logger *logrus.Logger

	pool       *pgxpool.Pool
	repository *Repository
}

func (a *app) parseUserID(r *http.Request) (*uuid.UUID, error) {
	strUserID := chi.URLParam(r, "id")

	if strUserID == "" {
		return nil, nil
	}

	userID, err := uuid.Parse(strUserID)
	if err != nil {
		a.logger.WithField("error", err).Debug(fmt.Sprintf("failed to parse userID (uuid) from: '%s'", strUserID))

		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("userID parsed: %s", userID))

	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.WithField("method", r.Method).Info("usersHandler called")

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
	a.logger.WithField("method", r.Method).Info("userHandler called")

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
	a.logger.WithField("method", r.Method).Info("userArticlesHandler called")

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

	a.logger.WithField("stacktrace", string(debug.Stack())).Panicf("panic!!!")
}

func (a *app) Init(ctx context.Context, logger *logrus.Logger) error {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse conn string (%s): %w", DatabaseURL, err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	config.ConnConfig.Logger = logrusadapter.NewLogger(logger) // логгер запросов в БД

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	a.logger = logger
	a.pool = pool
	a.repository = NewRepository(a.pool)

	return a.repository.InitSchema(ctx)
}

func (a *app) Serve() error {
	r := chi.NewRouter()

	r.Get("/users", http.HandlerFunc(a.usersHandler))
	r.Get("/users/{id}", http.HandlerFunc(a.userHandler))
	r.Get("/users/{id}/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/panic", http.HandlerFunc(a.panicHandler))

	return http.ListenAndServe("0.0.0.0:9000", r)
}
