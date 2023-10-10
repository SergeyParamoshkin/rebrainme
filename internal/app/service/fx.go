package service

import (
	"dumper/internal/app/httpsrv"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"service",
		fx.Provide(
			NewUser,
			NewServiceAdapter,
		),
		fx.Invoke(
			func(lc fx.Lifecycle) {
				lc.Append(fx.Hook{})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("service")
		}),
	)
}

type ServiceAdapter struct {
	fx.Out

	UserService httpsrv.UserService
}

func NewServiceAdapter(userService *User) ServiceAdapter {
	return ServiceAdapter{
		UserService: userService,
	}
}
