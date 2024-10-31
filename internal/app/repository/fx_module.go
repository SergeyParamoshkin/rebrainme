//nolint:ireturn // fx
package repository

import (
	"github.com/SergeyParamoshkin/alerts/internal/postgres"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	PG     *postgres.PG
}

func NewModule() fx.Option {
	return fx.Module(
		"",
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("alerts_repo")
		}),
	)
}
