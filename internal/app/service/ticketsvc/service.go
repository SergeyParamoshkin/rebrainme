package ticketsvc

import (
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	config *Config
	repo   Repo
}

func New(params Params) *Service {
	return &Service{
		logger: params.Logger,
		config: params.Config,
		repo:   params.Repo,
	}
}
