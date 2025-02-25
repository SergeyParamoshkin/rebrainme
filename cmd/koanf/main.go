package main

import (
	"fmt"
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Port int    `koanf:"port"`
	Host string `koanf:"host"`
}

func main() {
	k := koanf.New(".")

	// Загрузка конфигурации из JSON-файла
	if err := k.Load(file.Provider("config.json"), json.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Загрузка переменных окружения (переопределяют значения из файла)
	k.Load(env.Provider("APP_", ".", func(s string) string {
		return s
	}), nil)

	// Парсинг конфигурации в структуру
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	fmt.Printf("Config: %+v\n", config)
}
