package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DDL = `
		DROP TABLE IF EXISTS articles;
		DROP TABLE IF EXISTS users;
		
		CREATE TABLE IF NOT EXISTS users
		(
			id   uuid         NOT NULL
				CONSTRAINT users_pk
					PRIMARY KEY,
			name varchar(150) NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS articles
		(
			id      uuid         NOT NULL
				CONSTRAINT articles_pk
					PRIMARY KEY,
			title   varchar(150) NOT NULL,
			text    text         NOT NULL,
			user_id uuid
				CONSTRAINT articles_users_id_fk
					REFERENCES users
		);

		INSERT INTO public.users (id, name) VALUES ('b29f95a2-499a-4079-97f5-ff55c3854fcb', 'usr1');
		INSERT INTO public.users (id, name) VALUES ('b6dede74-ad09-4bb7-a036-997ab3ab3130', 'usr2');

		INSERT INTO public.articles (id, title, text, user_id) VALUES ('e4e12c87-88d8-413c-8ab6-57bfa4e953a8', 'article_11', 'some text', 'b29f95a2-499a-4079-97f5-ff55c3854fcb');
		INSERT INTO public.articles (id, title, text, user_id) VALUES ('68792339-715c-4823-a4d5-a85cefec8d36', 'article_12', 'hello, world!', 'b29f95a2-499a-4079-97f5-ff55c3854fcb');
		INSERT INTO public.articles (id, title, text, user_id) VALUES ('e095e3a2-5b8e-4bc8-b793-bc3606c4fdd5', 'article_21', 'why so serious?', 'b6dede74-ad09-4bb7-a036-997ab3ab3130');
	`

	UsersSelect        = `SELECT id, name FROM users`
	UserByIDSelect     = `SELECT id, name FROM users WHERE id = $1`
	UserArticlesSelect = `SELECT id, title, text, user_id FROM articles WHERE user_id = $1`
)

var (
	ErrNotFound      = errors.New("not found")
	ErrMultipleFound = errors.New("multiple found")
)

type Repository struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
	logger *slog.Logger
}

func (r *Repository) InitSchema(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Repository.InitSchema")
	defer span.End()

	r.logger.DebugContext(ctx, "Initializing database schema")

	_, err := r.pool.Exec(ctx, DDL)
	if err != nil {
		r.logger.ErrorContext(ctx, "Schema initialization failed", "error", err)
		return fmt.Errorf("schema init failed: %w", err)
	}

	r.logger.InfoContext(ctx, "Schema initialized successfully")
	return nil
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUser")
	defer span.End()

	logger := r.logger.With("userID", id)
	logger.DebugContext(ctx, "Executing user query", "query", UserByIDSelect)

	span.SetAttributes(
		attribute.String("query", UserByIDSelect),
		attribute.String("arg0", id.String()),
	)

	rows, err := r.pool.Query(ctx, UserByIDSelect, id)
	if err != nil {
		logger.ErrorContext(ctx, "Database query failed", "error", err)
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var (
		user  User
		found bool
	)

	for rows.Next() {
		if found {
			err := fmt.Errorf("%w: user id %s", ErrMultipleFound, id)
			logger.ErrorContext(ctx, "Multiple users found", "error", err)
			return nil, err
		}

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			logger.ErrorContext(ctx, "Row scan failed", "error", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		found = true
	}

	if err := rows.Err(); err != nil {
		logger.ErrorContext(ctx, "Rows iteration error", "error", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if !found {
		err := fmt.Errorf("%w: user id %s", ErrNotFound, id)
		logger.WarnContext(ctx, "User not found", "error", err)
		return nil, err
	}

	logger.DebugContext(ctx, "User found", "user", user)
	return &user, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]User, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUsers")
	defer span.End()

	r.logger.DebugContext(ctx, "Fetching all users", "query", UsersSelect)

	span.SetAttributes(attribute.String("query", UsersSelect))

	rows, err := r.pool.Query(ctx, UsersSelect)
	if err != nil {
		r.logger.ErrorContext(ctx, "Users query failed", "error", err)
		return nil, fmt.Errorf("users query failed: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			r.logger.ErrorContext(ctx, "User scan failed", "error", err)
			return nil, fmt.Errorf("user scan failed: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Rows iteration failed", "error", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	r.logger.DebugContext(ctx, "Users fetched", "count", len(users))
	return users, nil
}

func (r *Repository) GetUserArticles(ctx context.Context, userID uuid.UUID) ([]Article, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUserArticles")
	defer span.End()

	logger := r.logger.With("userID", userID)
	logger.DebugContext(ctx, "Fetching user articles", "query", UserArticlesSelect)

	span.SetAttributes(
		attribute.String("query", UserArticlesSelect),
		attribute.String("arg0", userID.String()),
	)

	rows, err := r.pool.Query(ctx, UserArticlesSelect, userID)
	if err != nil {
		logger.ErrorContext(ctx, "Articles query failed", "error", err)
		return nil, fmt.Errorf("articles query failed: %w", err)
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Text, &article.UserID); err != nil {
			logger.ErrorContext(ctx, "Article scan failed", "error", err)
			return nil, fmt.Errorf("article scan failed: %w", err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		logger.ErrorContext(ctx, "Rows iteration failed", "error", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	logger.DebugContext(ctx, "Articles fetched", "count", len(articles))
	return articles, nil
}

func NewRepository(pool *pgxpool.Pool, tracer trace.Tracer, logger *slog.Logger) *Repository {
	return &Repository{
		pool:   pool,
		tracer: tracer,
		logger: logger.With("component", "repository"),
	}
}
