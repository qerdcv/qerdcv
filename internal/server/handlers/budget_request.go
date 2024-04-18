package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateBudgetCategoryRequest struct {
	Name string `json:"name"`
}

func (r CreateBudgetCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
	)
}

type CreateBudgetTransactionRequest struct {
	CategoryID *string `json:"category_id"`
	Amount     float64 `json:"amount"`
}

func (r CreateBudgetTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CategoryID, validation.When(r.CategoryID != nil, is.UUIDv4)),
		validation.Field(&r.Amount, validation.Required, validation.Min(1.0), validation.Max(999999999.0)),
	)
}
