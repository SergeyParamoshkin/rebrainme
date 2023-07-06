package telegram

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"telegram",
		// Конструктор который потенциально вызывает для создания
		fx.Provide(
			NewTelegramConfig,
			NewTelegram,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, bot *Telegram) {
				lc.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						bot.StartBot()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						return nil
					},
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("mqtt")
		}),
	)
}
