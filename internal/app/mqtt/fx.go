package mqtt

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"mqtt",
		// Конструктор который потенциально вызывает для создания
		fx.Provide(
			NewConfig,
			NewMQTT,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, client *MQTT) {
				lc.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						return client.Start()
					},
					OnStop: func(ctx context.Context) error {
						return client.Stop()
					},
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("mqtt")
		}),
	)
}
