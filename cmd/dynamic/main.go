package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Port    int
	Host    string
	Logging Logging `yaml:"logging"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func main() {
	// Настройка Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Загружаем конфигурацию
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Динамическая перезагрузка
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		reloadConfig()
	})

	// Первоначальная загрузка конфигурации
	reloadConfig()

	// Бесконечный цикл для демонстрации
	for {
		time.Sleep(1 * time.Second)
	}
}

func reloadConfig() {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
	}

	fmt.Printf("Reloaded config: %+v\n", config)
}
