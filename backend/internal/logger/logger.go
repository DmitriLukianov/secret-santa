package logger

import (
	"log/slog"
	"os"
	"time"
)

const serviceName = "secret-santa"

func New(levelStr, stage string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLevel(levelStr),
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
			case slog.LevelKey:
				return slog.String("severity", a.Value.String())
			case slog.MessageKey:
				return slog.String("rest", a.Value.String())
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	log := slog.New(handler).With(
		slog.String("service", serviceName),
		slog.String("stage", stage),
	)

	return log
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
