package entities

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          string    `json:"id"`
	Supplier    string    `json:"supplier"`
	IssuedAt    time.Time `json:"issued_at"`
	Description string    `json:"description"`
	Qty         int       `json:"qty"`
	CostPerUnit float64   `json:"cost_per_unit"`
	Category    string    `json:"category"`
}

func (e *Expense) GenerateID() {
	if e.ID != "" {
		return
	}

	e.ID = uuid.NewString()
}

type ExpenseFilters struct {
	Pagination
	Text     *string
	Supplier *string
	From     *time.Time
	To       *time.Time
}

type PaginatedExpenses struct {
	Results []Expense
	Total   int64
}
