package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func RunServer(ctx context.Context, h http.Handler, logger *slog.Logger) error {
	s := http.Server{
		Addr:         "127.0.0.1:8080",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 2 * time.Second,
		Handler:      h,
	}

	go func() {
		logger.Info("running server", slog.String("addr", ":8080"))
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen and serve", slog.Any("error", err.Error()))
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("server close: %w", err)
	}

	return nil
}
