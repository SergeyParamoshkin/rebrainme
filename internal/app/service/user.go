package service

import (
	"context"
	"fmt"

	"dumper/internal/model"

	"github.com/google/uuid"
)

type User struct {
	userRepository UserRepository
}

func NewUser(userRepository UserRepository) *User {
	return &User{
		userRepository: userRepository,
	}
}

func (u *User) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error to get user by id: %w", err)
	}

	return user, nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	return fmt.Errorf("error creating user %w",
		u.userRepository.Create(ctx, user))
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	return fmt.Errorf("error updating user %w",
		u.userRepository.Update(ctx, user))
}

func (u *User) Delete(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("error deleting user %w",
		u.userRepository.Delete(ctx, id))
}
