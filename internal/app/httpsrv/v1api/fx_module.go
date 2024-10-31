//nolint:ireturn // fx
package v1api

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv"
	"github.com/SergeyParamoshkin/alerts/internal/app/service/ticketsvc"
	"github.com/SergeyParamoshkin/alerts/internal/tel"
)

type Params struct {
	fx.In

	TicketService TicketService

	Logger    *zap.Logger
	Telemetry *tel.Telemetry
	Config    *Config
}

type Result struct {
	fx.Out

	API httpsrv.API
}

func NewModule() fx.Option {
	return fx.Module(
		"api_v1",
		fx.Provide(
			New,
			newAdapter,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("v1api")
		}),
	)
}

type AdapterOut struct {
	fx.Out

	TicketService TicketService
}

func newAdapter(
	ts *ticketsvc.Service,
) AdapterOut {
	return AdapterOut{
		TicketService: ts,
	}
}
