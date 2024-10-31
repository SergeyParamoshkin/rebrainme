//nolint:ireturn // fx
package ticketsvc

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/SergeyParamoshkin/alerts/internal/app/repository"
	"github.com/SergeyParamoshkin/alerts/internal/postgres"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	PG     *postgres.PG
	Config *Config

	Repo Repo
}

func NewModule() fx.Option {
	return fx.Module(
		"ticket_service",
		fx.Provide(
			New,
			newAdapter,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("ticket_service")
		}),
	)
}

type AdapterOut struct {
	fx.Out

	Repo Repo
}

func newAdapter(repo *repository.Repo) AdapterOut {
	return AdapterOut{
		Repo: repo,
	}
}
