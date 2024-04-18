package handlers

type BudgetCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BudgetCategoryCreateResponse struct {
	Message string `json:"message"`
}

type BudgetCategoriesListResponse struct {
	Items []BudgetCategoryResponse `json:"items"`
}

type CreateBudgetTransactionResponse struct {
	Message string `json:"message"`
}
