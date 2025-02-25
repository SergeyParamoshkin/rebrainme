package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port int
	Host string
}

func main() {
	// Настройка Viper для работы с Consul
	viper.AddRemoteProvider("consul", "localhost:8500", "myapp/config")
	viper.SetConfigType("json")

	// Загружаем конфигурацию
	err := viper.ReadRemoteConfig()
	if err != nil {
		log.Fatalf("Error reading config from Consul: %s", err)
	}

	// Динамическая перезагрузка
	go func() {
		for {
			time.Sleep(5 * time.Second) // Проверяем изменения каждые 5 секунд
			err := viper.WatchRemoteConfig()
			if err != nil {
				log.Printf("Error watching remote config: %s", err)
			} else {
				reloadConfig()
			}
		}
	}()

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
