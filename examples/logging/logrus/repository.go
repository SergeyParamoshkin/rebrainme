package main

import (
	"context"
	"errors"
	"fmt"

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
	pool *pgxpool.Pool
}

func (r *Repository) InitSchema(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, DDL)
	return err
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	rows, _ := r.pool.Query(ctx, UserByIDSelect, id)

	var (
		user  User
		found bool
	)

	for rows.Next() {
		if found {
			return nil, fmt.Errorf("%w: user id %s", ErrMultipleFound, id)
		}

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}

		found = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("%w: user id %s", ErrNotFound, id)
	}

	return &user, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]User, error) {
	rows, _ := r.pool.Query(ctx, UsersSelect)

	ret := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}

		ret = append(ret, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Repository) GetUserArticles(ctx context.Context, userID uuid.UUID) ([]Article, error) {
	rows, _ := r.pool.Query(ctx, UserArticlesSelect, userID)

	ret := make([]Article, 0)

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.ID, &article.Title, &article.Text, &article.UserID); err != nil {
			return nil, err
		}

		ret = append(ret, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
