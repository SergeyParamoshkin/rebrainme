package main

import (
	"context"
	"log"

	"github.com/SergeyParamoshkin/rebrainme/examples/engine/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	app.Run(ctx)
}
