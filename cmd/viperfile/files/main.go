package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port int
	Host string
}

func main() {
	viper.SetConfigName("config") // Имя файла без расширения
	viper.SetConfigType("yaml")   // Формат файла
	viper.AddConfigPath(".")      // Путь к файлу

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	c := Config{
		Port: viper.GetInt("port"),
		Host: viper.GetString("host"),
	}

	fmt.Printf("Config: port=%d, host=%s\n", c.Port, c.Host)
}
