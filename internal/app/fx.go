package app

import (
	"dumper/internal/app/httpsrv"
	"dumper/internal/app/service"
	"dumper/internal/repository"
	"dumper/pkg/postgres"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			httpsrv.NewModule(),
			service.NewModule(),
			postgres.NewModule(),
			repository.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
