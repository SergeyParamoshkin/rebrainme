//nolint:ireturn // fx
package postgres

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Config *Config
}

func NewModule() fx.Option {
	return fx.Module(
		"postgres",
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("postgres")
		}),
		fx.Invoke(
			func(lc fx.Lifecycle, srv *PG) {
				lc.Append(fx.Hook{
					OnStart: srv.Start,
					OnStop:  srv.Stop,
				})
			},
		),
	)
}
