//nolint:ireturn // fx
package v1api

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
}

type Result struct {
	fx.Out

	API API
}

func NewModule() fx.Option {
	return fx.Module(
		"api_v1",
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("v1api")
		}),
	)
}
