package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Port int
	Host string
}

func main() {
	var config Config

	flag.IntVar(&config.Port, "port", 8080, "Port to run the server on")
	flag.StringVar(&config.Host, "host", "localhost", "Host to run the server on")
	flag.Parse()

	fmt.Printf("Config: %+v\n", config)
}
