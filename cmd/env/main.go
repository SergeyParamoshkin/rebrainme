package main

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
	Host string
}

func main() {
	config := Config{
		Port: os.Getenv("APP_PORT"),
		Host: os.Getenv("APP_HOST"),
	}

	fmt.Printf("Config: %+v\n", config)
}
