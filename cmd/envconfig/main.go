package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int    `envconfig:"PORT" default:"8080"`
	Host string `envconfig:"HOST" default:"localhost"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("Error processing env vars: %v", err)
	}

	fmt.Printf("Config: %+v\n", config)
}
