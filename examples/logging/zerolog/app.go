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
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type app struct {
	logger *zerolog.Logger

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
		a.logger.Err(err).Msg(
			fmt.Sprintf("failed to parse userID (uuid) from: '%s'", strUserID))

		return nil, err
	}

	a.logger.Debug().Str("msg", fmt.Sprintf("userID parsed: %s", userID)).Send()

	return &userID, nil
}

func (a *app) usersHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info().Interface("method", r.Method).Msg("usersHandler called")

	ctx := r.Context()

	users, err := a.repository.GetUsers(ctx)
	if err != nil {
		msg := fmt.Sprintf(`failed to get users: %s`, err)

		a.logger.Error().Err(err).Send()
		writeResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeJsonResponse(w, http.StatusOK, users)
}

func (a *app) userHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info().Interface("method", r.Method).Msg("usersHandler called")

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
	a.logger.Info().Interface("method", r.Method).Msg("userArticlesHandler called")

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

	a.logger.Panic().Msg(string(debug.Stack()))

}

func (a *app) Init(ctx context.Context, logger *zerolog.Logger) error {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse conn string (%s): %w", DatabaseURL, err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug

	config.ConnConfig.Logger = zerologadapter.NewLogger(*logger) // логгер запросов в БД

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
