package auth

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn // fx
	name := "auth_controller"
	return fx.Module(
		name,
		fx.Provide(
			NewAuthController,
			NewOAuth2Controller,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named(name)
		}),
	)
}
