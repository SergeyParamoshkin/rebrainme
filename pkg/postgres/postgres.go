package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Client struct {
	Logger *zap.Logger
	Config *Config
	Pool   *pgxpool.Pool
}

func NewPostgres(logger *zap.Logger, config *Config) (*Client, error) {
	return &Client{
		Logger: logger,
		Config: config,
		Pool:   nil,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Config.Host, c.Config.Port, c.Config.Username,
		c.Config.Password, c.Config.Database)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	c.Pool = pool
	err = c.Pool.Ping(ctx)

	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	return nil
}

func (c *Client) Stop(_ context.Context) error {
	c.Pool.Close()

	return nil
}
