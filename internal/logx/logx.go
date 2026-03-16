package logx

import (
	"fmt"
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

func ParseLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log level %q", value)
	}
}
