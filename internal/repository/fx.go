package repository

import (
	"dumper/internal/app/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option { //nolint:ireturn
	return fx.Module(
		"repository",

		fx.Provide(
			NewPostgresUserRepository,
			NewServiceAdapter,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("repository")
		}),
	)
}

type ServiceAdapter struct {
	fx.Out

	PostgresUserRepository service.UserRepository
}

func NewServiceAdapter(postgresUserRepository *PostgresUserRepository) ServiceAdapter {
	return ServiceAdapter{
		PostgresUserRepository: postgresUserRepository,
	}
}
