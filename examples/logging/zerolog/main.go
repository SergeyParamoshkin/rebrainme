package main

import (
	"context"

	"os"

	llog "log"

	"github.com/rs/zerolog"
)

const (
	DatabaseURL = "postgres://usr:pwd@localhost:5432/example?sslmode=disable"
)

func main() {

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	// logger.SetLevel(logrus.DebugLevel)
	// logger.SetFormatter(&logrus.JSONFormatter{})

	a := app{}

	if err := a.Init(context.Background(), &logger); err != nil {
		llog.Fatal(err)
	}

	if err := a.Serve(); err != nil {
		llog.Fatal(err)
	}
}
