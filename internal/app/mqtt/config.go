package mqtt

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Topic    string `yaml:"topic"`
	ClientID string `yaml:"clientID"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var mqttConfig Config
	err := provider.Get("mqtt").Populate(&mqttConfig)
	if err != nil {
		return nil, fmt.Errorf("mqtt config: %w", err)
	}

	return &mqttConfig, nil
}
