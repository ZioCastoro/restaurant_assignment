package payloads

import (
	"time"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
)

type ExpenseResponse struct {
	ID          string  `json:"id"`
	Supplier    string  `json:"supplier"`
	IssuedAt    string  `json:"issuedAt"`
	Description string  `json:"description"`
	Qty         int     `json:"qty"`
	CostPerUnit float64 `json:"costPerUnit"`
	Category    string  `json:"category"`
}

func NewExpenseResponse(entity entities.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:          entity.ID,
		Supplier:    entity.Supplier,
		IssuedAt:    entity.IssuedAt.Format(time.RFC3339),
		Description: entity.Description,
		Qty:         entity.Qty,
		CostPerUnit: entity.CostPerUnit,
		Category:    entity.Category,
	}
}

type PaginatedExpensesResponse struct {
	Results []ExpenseResponse `json:"results"`
	PaginationResponse
	Total int64 `json:"total"`
}

func NewPaginatedExpensesResponse(
	entity entities.PaginatedExpenses,
	pagination entities.Pagination,
) PaginatedExpensesResponse {
	results := make([]ExpenseResponse, 0, len(entity.Results))

	for _, expense := range entity.Results {
		results = append(results, NewExpenseResponse(expense))
	}

	return PaginatedExpensesResponse{
		Results:            results,
		PaginationResponse: NewPaginationResponse(pagination),
		Total:              entity.Total,
	}
}
