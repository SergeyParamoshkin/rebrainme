package main

import (
	"fmt"
	"os"

	"github.com/SergeyParamoshkin/rebrainme/internal/app"
	"github.com/SergeyParamoshkin/rebrainme/internal/config"
)

func main() {
	fmt.Println("Hello, RebrainMe!")

	c := config.NewConfig()
	if err := c.Validate(); err != nil {
		fmt.Printf("Error validating config: %v\n", err)
		os.Exit(1)
	}

	a := app.NewApp(c)
	err := a.Start()
	if err != nil {
		fmt.Printf("Error starting the app: %v\n", err)
		os.Exit(1)
	}
}
