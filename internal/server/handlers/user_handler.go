package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/qerdcv/qerdcv/internal/services"
)

type UserHandler struct {
	logger  *slog.Logger
	service *services.UserService
}

func NewUserHandler(
	logger *slog.Logger,
	service *services.UserService,
) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

func (h *UserHandler) AuthPage(c echo.Context) error {
	return c.Render(http.StatusOK, "templates/auth/index.gohtml", nil)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return EchoErrorFromValidation(err)
	}

	if err := h.service.CreateUser(
		c.Request().Context(),
		req.Username, req.Password,
	); err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, ErrorResponse{
				Message: services.ErrUserAlreadyExists.Error(),
			})
		}

		h.logger.Error("create user", slog.Any("err", err))
		return fmt.Errorf("service create user: %w", err)
	}

	return c.JSON(http.StatusCreated, RegisterResponse{Message: http.StatusText(http.StatusCreated)})
}

func (h *UserHandler) AuthorizeUser(c echo.Context) error {
	var req AuthorizeUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return EchoErrorFromValidation(err)
	}

	sessionToken, err := h.service.AuthorizeUser(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return fmt.Errorf("service authorize user: %w", err)
	}

	return c.JSON(http.StatusOK, AuthorizeUserResponse{
		SessionToken: sessionToken,
	})
}
