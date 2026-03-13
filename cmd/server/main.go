package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	serverapp "github.com/katsutoo/perissables/internal/app/server"
	"github.com/katsutoo/perissables/internal/logx"
)

const defaultServerAddr = ":8080"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := logx.New(logx.Config{
		Service: "server",
		Level:   logx.LevelFromEnv(os.Getenv("LOG_LEVEL")),
	})

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = defaultServerAddr
	}

	if err := serverapp.Run(ctx, logger, addr); err != nil {
		logger.Error("server exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}
