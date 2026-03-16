package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	clientapp "github.com/katsutoo/perissables/internal/app/client"
	"github.com/katsutoo/perissables/internal/config"
	"github.com/katsutoo/perissables/internal/logx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadClientFromEnv()
	if err != nil {
		slog.Error("load client configuration", slog.Any("err", err))
		os.Exit(1)
	}

	logger := logx.New(logx.Config{
		Service: "client",
		Level:   cfg.LogLevel,
	})

	if err := clientapp.Run(ctx, logger); err != nil {
		logger.Error("client exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}
