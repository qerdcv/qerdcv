package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/qerdcv/qerdcv/internal/services"
	"github.com/qerdcv/qerdcv/pkg/domain"
)

func Auth(userService *services.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if len(h) == 0 {
				return echo.ErrUnauthorized
			}

			parts := strings.Split(h, " ")
			if len(parts) != 2 {
				return echo.ErrUnauthorized
			}

			ctx := c.Request().Context()
			userSession, err := userService.VerifySession(c.Request().Context(), parts[1])
			if err != nil {
				if errors.Is(err, services.ErrInvalidToken) ||
					errors.Is(err, services.ErrSessionNotFound) ||
					errors.Is(err, services.ErrSessionExpired) {
					return echo.ErrUnauthorized
				}

				return fmt.Errorf("service verify session: %w", err)
			}

			c.SetRequest(c.Request().WithContext(domain.ContextWithUserSession(ctx, userSession)))

			return next(c)
		}
	}
}
