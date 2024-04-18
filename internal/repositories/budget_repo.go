package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/qerdcv/qerdcv/pkg/domain"
	"github.com/qerdcv/qerdcv/pkg/sqlutils"
)

type BudgetRepo struct {
	db sqlutils.DB
}

func NewBudgetRepo(db sqlutils.DB) *BudgetRepo {
	return &BudgetRepo{
		db: db,
	}
}

func (r *BudgetRepo) CreateCategory(ctx context.Context, uID int, name string) error {
	query := `INSERT INTO budget_categories(id, user_id, name) VALUES($1, $2, $3)`

	if _, err := r.db.ExecContext(ctx, query, uuid.New(), uID, name); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == sqlutils.PQErrCodeUniqueConstraint {
			return ErrUniqueConstraint
		}

		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}

func (r *BudgetRepo) CategoriesList(ctx context.Context, uID int) ([]domain.Category, error) {
	query := `SELECT id, name FROM budget_categories WHERE user_id=$1`

	rows, err := r.db.QueryContext(ctx, query, uID)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}

	defer rows.Close()

	categories := make([]domain.Category, 0)

	for rows.Next() {
		var category domain.Category
		if err = rows.Scan(
			&category.ID,
			&category.Name,
		); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		category.UserID = uID
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *BudgetRepo) CreateTransaction(ctx context.Context, uID int, cID *uuid.UUID, amount int64) error {
	query := `INSERT INTO budget_transactions (id, user_id, category_id, amount) VALUES ($1, $2, $3, $4)`

	if _, err := r.db.ExecContext(ctx, query, uuid.New(), uID, cID, amount); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == sqlutils.PQErrCodeForeignKeyConstraint {
			return ErrForeignKeyConstraint
		}

		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
