package httpsrv

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"http",
		// Конструктор который потенциально вызывает для создания
		fx.Provide(
			NewConfig,
			NewHTTPServer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, server *Server) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return server.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						return server.Stop(ctx)
					},
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("http")
		}),
	)
}
