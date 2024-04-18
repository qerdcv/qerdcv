package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						"panic recovered",
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack()[250:])),
					)
					err = &echo.HTTPError{
						Code:    http.StatusInternalServerError,
						Message: http.StatusText(http.StatusInternalServerError),
					}
				}
			}()

			err = next(c)
			return
		}
	}
}
