package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	clientapp "github.com/katsutoo/perissables/internal/app/client"
	"github.com/katsutoo/perissables/internal/logx"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := logx.New(logx.Config{
		Service: "client",
		Level:   logx.LevelFromEnv(os.Getenv("LOG_LEVEL")),
	})

	if err := clientapp.Run(ctx, logger); err != nil {
		logger.Error("client exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}
