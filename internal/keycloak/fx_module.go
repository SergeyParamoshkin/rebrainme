//nolint:ireturn // fx
package keycloak

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Config *Config
}

func NewModule() fx.Option {
	name := "keycloak"
	return fx.Module(
		name,
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("keycloak")
		}),
	)
}
