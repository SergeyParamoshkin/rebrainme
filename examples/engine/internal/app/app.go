package app

import (
	"context"
	"fmt"
	"net"

	"github.com/SergeyParamoshkin/rebrainme/examples/engine/internal/grpc/server"
	"github.com/SergeyParamoshkin/rebrainme/examples/engine/internal/logger"
	"go.uber.org/zap"
)

type App struct {
	Log logger.Logger
}

func NewApp(ctx context.Context) (*App, error) {
	zapLogger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &App{
		Log: logger.NewZapLogger(zapLogger),
	}, nil
}

func (a *App) Run(ctx context.Context) {
	a.Log.Info(fmt.Errorf("Starting listening on port 8080"))
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		a.Log.Info(fmt.Errorf("failed to listen: %w", err))
	}

	a.Log.Info(fmt.Errorf("Listening on %s", port))

	srv := server.NewServer()

	if err := srv.Serve(lis); err != nil {
		a.Log.Error(fmt.Errorf("failed to serve: %w", err))
	}
}
