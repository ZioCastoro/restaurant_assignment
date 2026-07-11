package payloads

import (
	"time"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/app/validators"
)

type ExpenseRequest struct {
	Supplier    string  `json:"supplier,omitempty" validate:"required,notblank"`
	IssuedAt    string  `json:"issuedAt,omitempty" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Description string  `json:"description,omitempty"`
	Qty         int     `json:"qty,omitempty" validate:"required,gt=0"`
	CostPerUnit float64 `json:"costPerUnit,omitempty" validate:"required,gt=0,cent"`
	Category    string  `json:"category,omitempty" validate:"required,notblank"`
}

func (e *ExpenseRequest) Validate() error {
	return validators.ValidateStruct(e)
}

func (e *ExpenseRequest) MapToEntity() entities.Expense {
	// error is ignored because this method is intended to only be used with validated data.
	issuedAt, _ := time.Parse(time.RFC3339, e.IssuedAt)

	return entities.Expense{
		Supplier:    e.Supplier,
		IssuedAt:    issuedAt,
		Description: e.Description,
		Qty:         e.Qty,
		CostPerUnit: e.CostPerUnit,
		Category:    e.Category,
	}
}

type CreateExpenseRequest struct {
	ExpenseRequest
}

type UpdateExpenseRequest struct {
	ExpenseRequest
}

func (e *UpdateExpenseRequest) MapToEntity(id string) entities.Expense {
	entity := e.ExpenseRequest.MapToEntity()
	entity.ID = id

	return entity
}

type ExpenseFilterRequest struct {
	PaginationRequest
	Text     *string `query:"text,omitempty" validate:"omitempty,notblank"`
	Supplier *string `query:"supplier,omitempty" validate:"omitempty,notblank"`
	From     *string `query:"from,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	To       *string `query:"to,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (e *ExpenseFilterRequest) Validate() error {
	return validators.ValidateStruct(e)
}

func (e *ExpenseFilterRequest) MapToEntity() entities.ExpenseFilters {
	return entities.ExpenseFilters{
		Pagination: e.PaginationRequest.MapToEntity(),
		Text:       e.Text,
		Supplier:   e.Supplier,
		From:       parseTime(e.From),
		To:         parseTime(e.To),
	}
}

type ImportExpenseRequest []CreateExpenseRequest

func (e *ImportExpenseRequest) Validate() error {
	return validators.ValidateSlice(e)
}

func (e *ImportExpenseRequest) MapToEntity() []entities.Expense {
	expenses := make([]entities.Expense, 0, len(*e))

	for _, expense := range *e {
		expenses = append(expenses, expense.MapToEntity())
	}

	return expenses
}
