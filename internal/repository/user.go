package repository

import (
	"context"
	"fmt"

	"dumper/internal/model"

	"dumper/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

type PostgresUserRepository struct {
	Client *postgres.Client
}

func NewPostgresUserRepository(client *postgres.Client) *PostgresUserRepository {
	return &PostgresUserRepository{
		Client: client,
	}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query, _, err := squirrel.Select("*").
		From("User").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.Client.Pool.QueryRow(ctx, query)

	var User model.User
	err = row.Scan(&User.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to scan User: %w", err)
	}

	return &User, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, User *model.User) error {
	query, args, err := squirrel.Insert("User").
		Columns("id").
		Values(User.ID).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.Client.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create User: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, User *model.User) error {
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
