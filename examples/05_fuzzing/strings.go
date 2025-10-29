package fuzzing

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

func Reverse(s string) string {
	if !utf8.ValidString(s) {
		return s
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func IsPalindrome(s string) bool {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, strings.ToLower(s))

	runes := []rune(cleaned)
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

func ParseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Port < 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", cfg.Port)
	}

	return &cfg, nil
}

func ParseURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	return u, nil
}

func SafeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username too short")
	}
	if len(username) > 20 {
		return fmt.Errorf("username too long")
	}
	if !utf8.ValidString(username) {
		return fmt.Errorf("invalid UTF-8")
	}

	for _, r := range username {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return fmt.Errorf("invalid character: %c", r)
		}
	}

	return nil
}

func ParseKeyValue(input string) (map[string]string, error) {
	result := make(map[string]string)

	if input == "" {
		return result, nil
	}

	pairs := strings.Split(input, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key in pair: %s", pair)
		}

		result[key] = value
	}

	return result, nil
}
