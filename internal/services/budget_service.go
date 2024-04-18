package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/qerdcv/qerdcv/internal/repositories"
	"github.com/qerdcv/qerdcv/pkg/domain"
)

var (
	ErrBudgetCategoryAlreadyExists = errors.New("budget category already exists")
	ErrBudgetCategoryNotFound      = errors.New("budget category not found")
)

type BudgetService struct {
	repo *repositories.BudgetRepo
}

func NewBudgetService(repo *repositories.BudgetRepo) *BudgetService {
	return &BudgetService{
		repo: repo,
	}
}

func (s *BudgetService) CreateCategory(ctx context.Context, uID int, name string) error {
	if err := s.repo.CreateCategory(ctx, uID, name); err != nil {
		if errors.Is(err, repositories.ErrUniqueConstraint) {
			return ErrBudgetCategoryAlreadyExists
		}

		return fmt.Errorf("repo create category: %w", err)
	}

	return nil
}

func (s *BudgetService) CategoriesList(ctx context.Context, uID int) ([]domain.Category, error) {
	categories, err := s.repo.CategoriesList(ctx, uID)
	if err != nil {
		return nil, fmt.Errorf("repo categories list: %w", err)
	}

	return categories, nil
}

func (s *BudgetService) CreateTransaction(ctx context.Context, uID int, cID *uuid.UUID, amount int64) error {
	if err := s.repo.CreateTransaction(ctx, uID, cID, amount); err != nil {
		if errors.Is(err, repositories.ErrForeignKeyConstraint) {
			return ErrBudgetCategoryNotFound
		}

		return fmt.Errorf("repo create transaction: %w", err)
	}

	return nil
}
