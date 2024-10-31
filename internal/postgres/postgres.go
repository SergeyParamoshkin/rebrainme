package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PG struct {
	*pgxpool.Pool

	config *Config
}

func New(params Params) (*PG, error) {
	return &PG{
		Pool:   nil,
		config: params.Config,
	}, nil
}

func (p *PG) Start(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(p.config.DatabaseURL)
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return err
	}

	p.Pool = pool

	return nil
}

func (p *PG) Stop(ctx context.Context) error {
	p.Pool.Close()

	return nil
}
