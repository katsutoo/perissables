package logx

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Service string
	Level   slog.Level
	Writer  io.Writer
}

func New(config Config) *slog.Logger {
	if config.Service == "" {
		config.Service = "app"
	}

	if config.Writer == nil {
		config.Writer = os.Stderr
	}

	handler := slog.NewTextHandler(config.Writer, &slog.HandlerOptions{Level: config.Level})

	return slog.New(handler).With(slog.String("service", config.Service))
}

func LevelFromEnv(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
