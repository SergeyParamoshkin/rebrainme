package telegram

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Token  string `yaml:"token"`
	Debug  bool   `yaml:"debug"`
	ChanID int64  `yaml:"chatID"`
}

func NewTelegramConfig(provider config.Provider) (*Config, error) {
	var config Config
	err := provider.Get("telegram").Populate(&config)
	if err != nil {
		return nil, fmt.Errorf("telegram config: %w", err)
	}

	return &config, nil
}
