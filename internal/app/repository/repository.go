package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	fwctx "github.com/SergeyParamoshkin/alerts/internal/ctx"
	"github.com/SergeyParamoshkin/alerts/internal/postgres"
)

type Repo struct {
	logger *zap.Logger
	pg     *postgres.PG

	stmtBuilder squirrel.StatementBuilderType
}

func New(params Params) *Repo {
	stmtBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Repo{
		logger:      params.Logger,
		pg:          params.PG,
		stmtBuilder: stmtBuilder,
	}
}

func (r *Repo) commitRollback(ctx context.Context, tx pgx.Tx, err error) error {
	if err == nil {
		return tx.Commit(ctx)
	}

	if rErr := tx.Rollback(ctx); rErr != nil {
		fwctx.LoggerFromCtx(ctx, r.logger).Error("failed to rollback transaction", zap.Error(rErr))
	}

	return err
}

func (r *Repo) WithConn(
	ctx context.Context,
	callback func(ctx context.Context, q Queryable) error,
) error {
	poolConn, err := r.pg.Acquire(ctx)
	if err != nil {
		return err
	}

	defer poolConn.Release()

	conn := poolConn.Conn()

	return callback(ctx, conn)
}

func (r *Repo) WithTx(
	ctx context.Context,
	callback func(ctx context.Context, q Queryable) error,
) error {
	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() { err = r.commitRollback(ctx, tx, err) }()

	return callback(ctx, tx)
}
