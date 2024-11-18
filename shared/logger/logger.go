package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func getLogLevel(logLevel string) slog.Leveler {
	var level slog.Leveler
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return level
}

func New(logLevel string) *slog.Logger {
	// enable pretty-print logs
	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339Nano, // Full timestamp with nanoseconds
	}

	zerologLogger := zerolog.New(writer).Level(mapLevel(getLogLevel(logLevel).Level()))
	zerologLogger = zerologLogger.With().Timestamp().Logger()

	logger := slog.New(newZerologHandler(&zerologLogger))

	return logger
}
