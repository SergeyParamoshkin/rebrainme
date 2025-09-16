package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

var Logger *slog.Logger

func Init(debug bool) {
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Форматируем время в читаемом виде
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(time.Now().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	}

	// Используем JSON handler для production, Text handler для разработки
	var handler slog.Handler
	if debug {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)

	Logger.Info("Logger initialized",
		slog.String("level", level.String()),
		slog.Bool("debug", debug),
	)
}

// Helper функции для удобного логирования

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, err error, args ...any) {
	allArgs := append([]any{slog.String("error", err.Error())}, args...)
	Logger.Error(msg, allArgs...)
}

func Fatal(msg string, err error, args ...any) {
	allArgs := append([]any{slog.String("error", err.Error())}, args...)
	Logger.Error(msg, allArgs...)
	os.Exit(1)
}

// WithContext добавляет контекстную информацию к логгеру
func WithContext(ctx context.Context) *slog.Logger {
	// Извлекаем request ID и другие метаданные из контекста
	if reqID := ctx.Value("requestID"); reqID != nil {
		return Logger.With(slog.String("request_id", reqID.(string)))
	}
	return Logger
}

// LogRequest логирует информацию о запросе
func LogRequest(method, path string, statusCode int, duration time.Duration, args ...any) {
	baseArgs := []any{
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", statusCode),
		slog.Duration("duration", duration),
	}
	allArgs := append(baseArgs, args...)

	if statusCode >= 500 {
		Logger.Error("Request completed with error", allArgs...)
	} else if statusCode >= 400 {
		Logger.Warn("Request completed with client error", allArgs...)
	} else {
		Logger.Info("Request completed", allArgs...)
	}
}