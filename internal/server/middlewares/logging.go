package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Logging(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			level := slog.LevelInfo
			start := time.Now()

			err := next(c)

			attrs := []any{
				slog.String("path", c.Request().URL.Path),
				slog.String("method", c.Request().Method),
				slog.Duration("duration", time.Since(start)),
			}
			if err != nil {
				var echoErr *echo.HTTPError
				if errors.As(err, &echoErr) {
					if echoErr.Code >= http.StatusInternalServerError {
						level = slog.LevelError
						attrs = append(attrs, slog.Any("err", err))
					}
					attrs = append(attrs, slog.Int("status_code", echoErr.Code))
				} else {
					attrs = append(attrs, slog.Any("err", err))
				}
			} else {
				attrs = append(attrs, slog.Int("status_code", c.Response().Status))
			}

			logger.Log(c.Request().Context(), level, "handled request", attrs...)

			return err
		}
	}
}
