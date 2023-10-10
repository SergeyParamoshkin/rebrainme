package service

import (
	"context"

	"dumper/internal/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, report *model.User) error
	Update(ctx context.Context, report *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
