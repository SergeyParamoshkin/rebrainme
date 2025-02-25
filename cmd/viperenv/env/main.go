package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	viper.AutomaticEnv() // Переменные окружения переопределяют значения из файла

	port := viper.GetInt("port")
	host := viper.GetString("host")

	fmt.Printf("Config: port=%d, host=%s\n", port, host)
}
