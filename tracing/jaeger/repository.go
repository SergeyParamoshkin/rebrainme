package main

import (
	"context"
	"errors"
	"fmt"

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

var ErrUserNotFound = errors.New("User not found")

// Custom error types
type NotFoundError struct {
	Entity string
	ID     uuid.UUID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Entity, e.ID)
}

type MultipleFoundError struct {
	Entity string
	ID     uuid.UUID
}

func (e *MultipleFoundError) Error() string {
	return fmt.Sprintf("multiple %s found with id %s", e.Entity, e.ID)
}

// Для поддержки errors.Is
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

type RepositoryError struct {
	Operation string
	Err       error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("repository error during %s: %v", e.Operation, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

type Repository struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
}

func (r *Repository) InitSchema(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Repository.InitSchema")
	defer span.End()

	_, err := r.pool.Exec(ctx, DDL)
	if err != nil {
		return &RepositoryError{
			Operation: "schema initialization",
			Err:       err,
		}
	}
	return nil
}

// Пример использования с проверкой типов ошибок
func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := r.getUser(ctx, id)
	if err != nil {
		switch e := err.(type) {
		case *NotFoundError:
			return nil, &NotFoundError{}
			// Специфичная обработка для NotFound
		case *MultipleFoundError:
			fmt.Printf("Multiple found: %v\n", e)
			// Специфичная обработка для MultipleFound
		case *RepositoryError:
			fmt.Printf("Repository error: %v\n", e)
			return nil, &RepositoryError{}
			// Специфичная обработка для RepositoryError
		default:
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
	return user, nil
}

func (r *Repository) getUser(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("query", UserByIDSelect),
		attribute.String("arg0", id.String()),
	)

	rows, err := r.pool.Query(ctx, UserByIDSelect, id)
	if err != nil {
		return nil, &RepositoryError{
			Operation: "get user query",
			Err:       err,
		}
	}

	defer rows.Close()

	var (
		user  User
		found bool
	)

	for rows.Next() {
		if found {
			return nil, &MultipleFoundError{
				Entity: "user",
				ID:     id,
			}
		}

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, &RepositoryError{
				Operation: "user row scan",
				Err:       err,
			}
		}

		found = true
	}

	if err := rows.Err(); err != nil {
		return nil, &RepositoryError{
			Operation: "rows processing",
			Err:       err,
		}
	}

	if !found {
		// return nil, &NotFoundError{
		// 	Entity: "user",
		// 	ID:     id,
		// }
		return nil, fmt.Errorf("user %s not found %w", id, ErrUserNotFound)
	}

	return &user, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]User, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUsers")
	defer span.End()

	span.SetAttributes(
		attribute.String("query", UsersSelect),
	)

	rows, err := r.pool.Query(ctx, UsersSelect)
	if err != nil {
		return nil, &RepositoryError{
			Operation: "get users query",
			Err:       err,
		}
	}
	defer rows.Close()

	ret := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, &RepositoryError{
				Operation: "users row scan",
				Err:       err,
			}
		}

		ret = append(ret, user)
	}

	if err := rows.Err(); err != nil {
		return nil, &RepositoryError{
			Operation: "rows processing",
			Err:       err,
		}
	}

	return ret, nil
}

func (r *Repository) GetUserArticles(ctx context.Context, userID uuid.UUID) ([]Article, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUserArticles")
	defer span.End()

	span.SetAttributes(
		attribute.String("query", UserArticlesSelect),
		attribute.String("arg0", userID.String()),
	)

	rows, err := r.pool.Query(ctx, UserArticlesSelect, userID)
	if err != nil {
		return nil, &RepositoryError{
			Operation: "get user articles query",
			Err:       err,
		}
	}
	defer rows.Close()

	ret := make([]Article, 0)

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.ID, &article.Title, &article.Text, &article.UserID); err != nil {
			return nil, &RepositoryError{
				Operation: "article row scan",
				Err:       err,
			}
		}

		ret = append(ret, article)
	}

	if err := rows.Err(); err != nil {
		return nil, &RepositoryError{
			Operation: "rows processing",
			Err:       err,
		}
	}

	return ret, nil
}

func NewRepository(pool *pgxpool.Pool, tracer trace.Tracer) *Repository {
	return &Repository{pool: pool, tracer: tracer}
}
