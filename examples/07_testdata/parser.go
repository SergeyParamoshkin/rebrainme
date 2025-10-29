package testdata

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout"`
	Debug   bool   `json:"debug"`
}

func ParseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}

	return nil
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func ParseUsers(data []byte) ([]User, error) {
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	for i, user := range users {
		if user.Name == "" {
			return nil, fmt.Errorf("user at index %d has empty name", i)
		}
		if user.Email == "" {
			return nil, fmt.Errorf("user at index %d has empty email", i)
		}
	}

	return users, nil
}

type Report struct {
	Title   string
	Date    string
	Summary string
	Items   []string
}

func GenerateReport(r *Report) string {
	result := fmt.Sprintf("=== %s ===\n", r.Title)
	result += fmt.Sprintf("Date: %s\n\n", r.Date)
	result += fmt.Sprintf("%s\n\n", r.Summary)

	if len(r.Items) > 0 {
		result += "Items:\n"
		for i, item := range r.Items {
			result += fmt.Sprintf("  %d. %s\n", i+1, item)
		}
	}

	return result
}
