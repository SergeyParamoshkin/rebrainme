package app

import (
	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv/v1api"
	"github.com/SergeyParamoshkin/alerts/internal/app/repository"
	"github.com/SergeyParamoshkin/alerts/internal/app/service/ticketsvc"
	"github.com/SergeyParamoshkin/alerts/internal/config"
	"github.com/SergeyParamoshkin/alerts/internal/postgres"
	"github.com/SergeyParamoshkin/alerts/internal/tel"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewApp(config config.Config) *fx.App {
	return fx.New(
		httpsrv.NewModule(),
		v1api.NewModule(),
		ticketsvc.NewModule(),
		repository.NewModule(),
		postgres.NewModule(),
		fx.Provide(
			func(config *AppConfig) (*zap.Logger, error) {
				if config.Debug {
					return zap.NewDevelopment()
				}

				return zap.NewProduction()
			},
			NewArgs,
			NewConfig,
			tel.New,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
