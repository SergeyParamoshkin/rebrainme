package app

import (
	"dumper/internal/app/httpsrv"
	"dumper/internal/app/mqtt"
	"dumper/internal/app/telegram"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			mqtt.NewModule(),
			telegram.NewModule(),
			httpsrv.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig, // Провайдер для загрузки конфигурации
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
