package repositories

import (
	"context"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
)

type ExpensesRepositoryInterface interface {
	List(ctx context.Context, filters entities.ExpenseFilters) (entities.PaginatedExpenses, error)
	Find(ctx context.Context, id string) (entities.Expense, error)
	Create(ctx context.Context, expense entities.Expense) (entities.Expense, error)
	Update(ctx context.Context, expense entities.Expense) (entities.Expense, error)
	Delete(ctx context.Context, id string) error
}
