package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

func main() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	fmt.Printf("Config: %+v\n", config)
}
