package http

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			h.ServeHTTP(w, r)
			logger.Info(
				"handling request",
				slog.String("url", r.URL.String()),
				slog.String("ip", r.RemoteAddr),
				slog.Duration("duration", time.Since(start)),
				slog.Int64("content_length", r.ContentLength),
			)
		})
	}
}
