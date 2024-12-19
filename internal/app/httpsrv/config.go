package httpsrv

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/config"
)

type Config struct {
	Port              string        `yaml:"port"`
	Host              string        `yaml:"host"`
	MetricPath        string        `yaml:"metricPath"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

func (c *Config) Validate() error {
	context.TODO()

	if c.Port == "" {
		return errors.New("port is required")
	}
	if _, err := net.LookupPort("tcp", c.Port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}

	if c.Host == "" {
		return errors.New("host is required")
	}

	if c.MetricPath == "" {
		return errors.New("metricPath is required")
	}

	if c.ReadHeaderTimeout <= 0 {
		return errors.New("readHeaderTimeout must be greater than zero")
	}

	return nil
}

func NewConfig(provider config.Provider) (*Config, error) {
	var config Config
	err := provider.Get("http").Populate(&config)
	if err != nil {
		return &Config{}, fmt.Errorf("provider error: %w", err)
	}

	if err := config.Validate(); err != nil {
		return &Config{}, fmt.Errorf("validation error: %w", err)
	}

	return &config, nil
}
