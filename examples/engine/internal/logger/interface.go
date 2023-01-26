package logger

import "context"

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Logger interface {
	Debug(err error, tags ...Tag)
	Info(err error, tags ...Tag)
	Error(err error, tags ...Tag)

	With(tags ...Tag) Logger
	WithContext(ctx context.Context) Logger

	NewContext(ctx context.Context, tags ...Tag) context.Context
}
