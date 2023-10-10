package postgres

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

func NewPostgresConfig(provider config.Provider) (*Config, error) {
	var config Config
	err := provider.Get("postgres").Populate(&config)
	if err != nil {
		return nil, fmt.Errorf("postgres config: %w", err)
	}

	return &config, nil
}
