package client

import (
	"context"
	"log/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	logger.Info("client bootstrap ready")
	<-ctx.Done()
	logger.Info("client shutdown requested")

	return nil
}
