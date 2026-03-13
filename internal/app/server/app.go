package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const shutdownTimeout = 10 * time.Second

type healthResponse struct {
	Status string `json:"status"`
}

func Run(ctx context.Context, logger *slog.Logger, addr string) error {
	if logger == nil {
		logger = slog.Default()
	}

	if addr == "" {
		return errors.New("server address is required")
	}

	router := chi.NewRouter()
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
			logger.Error("write health response", slog.Any("err", err))
		}
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("server listening", slog.String("addr", addr))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen server: %w", err)
			return
		}

		close(errCh)
	}()

	select {
	case <-ctx.Done():
		logger.Info("server shutdown requested")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		if err, ok := <-errCh; ok {
			return err
		}

		return nil
	case err, ok := <-errCh:
		if !ok {
			return nil
		}

		return err
	}
}
