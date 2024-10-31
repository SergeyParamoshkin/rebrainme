//nolint:ireturn // fx
package httpsrv

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Config *Config
	API    API
}

func NewModule() fx.Option {
	return fx.Module(
		"httpsrv",
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("http_server")
		}),
		fx.Invoke(
			func(lc fx.Lifecycle, srv *Server) {
				lc.Append(fx.Hook{
					OnStart: srv.Start,
					OnStop:  srv.Stop,
				})
			},
		),
	)
}
