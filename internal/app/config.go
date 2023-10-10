package app

import (
	"flag"
	"fmt"

	"go.uber.org/config"
	"go.uber.org/fx"
)

type Config struct {
	Name string `yaml:"name"`
}

type ResultConfig struct {
	fx.Out

	Provider config.Provider
	Config   Config
}

func NewConfig() (ResultConfig, error) {
	// TODO: rewrite this to use a config file
	configFileName := ""
	flag.StringVar(&configFileName, "config", "config.yml", "config file")
	flag.Parse()

	loader, err := config.NewYAML(config.File(configFileName))
	if err != nil {
		return ResultConfig{}, fmt.Errorf("failed to load config file: %w", err)
	}

	config := Config{
		Name: "default",
	}

	if err := loader.Get("app").Populate(&config); err != nil {
		return ResultConfig{}, fmt.Errorf("failed to populate config: %w", err)
	}

	return ResultConfig{
		Provider: loader,
		Config:   config,
	}, nil
}
