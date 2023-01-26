package main

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	InitSchema(ctx context.Context) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserArticles(ctx context.Context, userID uuid.UUID) ([]Article, error)
}
