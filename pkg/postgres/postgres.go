package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Client struct {
	Logger *zap.Logger
	Pool   *pgxpool.Pool
}

func NewPostgres(logger *zap.Logger) (*Client, error) {
	return &Client{
		Logger: logger,
		Pool:   nil,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "postgres", "postgres")

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}

	c.Pool = pool

	return nil
}

func (c *Client) Stop(_ context.Context) error {
	c.Pool.Close()

	return nil
}
