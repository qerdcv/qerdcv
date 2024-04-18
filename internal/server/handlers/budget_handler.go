package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"

	"github.com/qerdcv/qerdcv/internal/services"
	"github.com/qerdcv/qerdcv/pkg/domain"
)

var decMul = decimal.NewFromFloat(100.00)

type BudgetHandler struct {
	logger  *slog.Logger
	service *services.BudgetService
}

func NewBudgetHandler(
	logger *slog.Logger,
	service *services.BudgetService,
) *BudgetHandler {
	return &BudgetHandler{
		logger:  logger,
		service: service,
	}
}

func (h *BudgetHandler) CreateCategory(c echo.Context) error {
	ctx := c.Request().Context()
	var req CreateBudgetCategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return EchoErrorFromValidation(err)
	}

	s := domain.UserSessionFromContext(ctx)

	if err := h.service.CreateCategory(ctx, s.UserID, req.Name); err != nil {
		if errors.Is(err, services.ErrBudgetCategoryAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, ErrorResponse{
				Message: services.ErrBudgetCategoryAlreadyExists.Error(),
			})
		}

		return fmt.Errorf("service create category: %w", err)
	}

	return c.JSON(http.StatusCreated, BudgetCategoryCreateResponse{
		Message: http.StatusText(http.StatusCreated),
	})
}

func (h *BudgetHandler) CategoriesList(c echo.Context) error {
	ctx := c.Request().Context()
	session := domain.UserSessionFromContext(ctx)

	categories, err := h.service.CategoriesList(c.Request().Context(), session.UserID)
	if err != nil {
		return fmt.Errorf("service categories list: %w", err)
	}

	respCategories := make([]BudgetCategoryResponse, len(categories))
	for i, ctg := range categories {
		respCategories[i] = BudgetCategoryResponse{
			ID:   ctg.ID.String(),
			Name: ctg.Name,
		}
	}

	return c.JSON(http.StatusOK, BudgetCategoriesListResponse{
		Items: respCategories,
	})
}

func (h *BudgetHandler) CreateTransaction(c echo.Context) error {
	var req CreateBudgetTransactionRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return EchoErrorFromValidation(err)
	}

	var cID *uuid.UUID
	if req.CategoryID != nil {
		id := uuid.MustParse(*req.CategoryID)
		cID = &id
	}

	ctx := c.Request().Context()
	s := domain.UserSessionFromContext(ctx)
	if err := h.service.CreateTransaction(
		ctx,
		s.UserID, cID, decimal.NewFromFloat(req.Amount).Mul(decMul).IntPart(),
	); err != nil {
		if errors.Is(err, services.ErrBudgetCategoryNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, services.ErrBudgetCategoryNotFound.Error())
		}

		return fmt.Errorf("service create transaction: %w", err)
	}

	return c.JSON(http.StatusCreated, CreateBudgetTransactionResponse{
		Message: http.StatusText(http.StatusCreated),
	})
}

func (h *BudgetHandler) TransactionsList(c echo.Context) error {
	return nil
}
