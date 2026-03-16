package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/katsutoo/perissables/internal/logx"
)

const defaultServerAddr = ":8080"

type Server struct {
	LogLevel slog.Level
	Addr     string
}

type Client struct {
	LogLevel slog.Level
}

func LoadServerFromEnv() (Server, error) {
	logLevel, err := parseLogLevelFromEnv()
	if err != nil {
		return Server{}, err
	}

	addr := strings.TrimSpace(os.Getenv("SERVER_ADDR"))
	if addr == "" {
		addr = defaultServerAddr
	}

	return Server{
		LogLevel: logLevel,
		Addr:     addr,
	}, nil
}

func LoadClientFromEnv() (Client, error) {
	logLevel, err := parseLogLevelFromEnv()
	if err != nil {
		return Client{}, err
	}

	return Client{LogLevel: logLevel}, nil
}

func parseLogLevelFromEnv() (slog.Level, error) {
	level, err := logx.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return 0, fmt.Errorf("parse LOG_LEVEL: %w", err)
	}

	return level, nil
}
