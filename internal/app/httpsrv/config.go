package httpsrv

import (
	"fmt"
	"time"

	"go.uber.org/config"
)

type Config struct {
	Port              string        `yaml:"port"`
	Host              string        `yaml:"host"`
	MetricPath        string        `yaml:"metricPath"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var config Config
	err := provider.Get("http").Populate(&config)
	if err != nil {
		return &Config{}, fmt.Errorf("provider error: %w", err)
	}

	return &config, nil
}
