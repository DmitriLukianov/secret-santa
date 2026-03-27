package logger

import (
	"log/slog"
	"os"
	"runtime"
	"time"
)

const (
	serviceName = "secret-santa"
)

func New(levelStr, stage string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: parseLevel(levelStr),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   "timestamp",
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			if a.Key == slog.LevelKey {
				return slog.Attr{
					Key:   "severity",
					Value: a.Value,
				}
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
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func WithSource(log *slog.Logger) *slog.Logger {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return log
	}
	fn := runtime.FuncForPC(pc)

	return log.With(
		slog.Group("source",
			slog.String("function", fn.Name()),
			slog.String("file", file),
			slog.Int("line", line),
		),
	)
}
