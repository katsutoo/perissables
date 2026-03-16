package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	serverapp "github.com/katsutoo/perissables/internal/app/server"
	"github.com/katsutoo/perissables/internal/config"
	"github.com/katsutoo/perissables/internal/logx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadServerFromEnv()
	if err != nil {
		slog.Error("load server configuration", slog.Any("err", err))
		os.Exit(1)
	}

	logger := logx.New(logx.Config{
		Service: "server",
		Level:   cfg.LogLevel,
	})

	if err := serverapp.Run(ctx, logger, cfg.Addr); err != nil {
		logger.Error("server exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}
