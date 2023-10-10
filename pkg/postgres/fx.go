package postgres

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"postgres",

		fx.Provide(
			NewPostgresConfig,
			NewPostgres,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, postgres *Client) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return postgres.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						return postgres.Stop(ctx)
					},
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("postgres")
		}))
}
