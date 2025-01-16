package app

import "github.com/SergeyParamoshkin/rebrainme/internal/config"

type App struct {
	Config config.Config
}

func NewApp(config config.Config) *App {
	return &App{
		Config: config,
	}
}

func (a *App) Start() error {
	return nil
}
