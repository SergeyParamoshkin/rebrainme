package main

import (
	"github.com/SergeyParamoshkin/alerts/internal/app"
	"github.com/SergeyParamoshkin/alerts/internal/config"
)

func main() {
	c := config.Config{}
	a := app.NewApp(c)
	a.Run()
}
