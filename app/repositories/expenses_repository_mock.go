package repositories

import (
	"context"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/stretchr/testify/mock"
)

type ExpensesRepositoryMock struct {
	ExpensesRepositoryInterface
	mock.Mock
}

func (m *ExpensesRepositoryMock) List(
	ctx context.Context,
	filters entities.ExpenseFilters,
) (entities.PaginatedExpenses, error) {
	args := m.Called(ctx, filters)

	return args.Get(0).(entities.PaginatedExpenses), args.Error(1)
}

func (m *ExpensesRepositoryMock) Find(
	ctx context.Context,
	id string,
) (entities.Expense, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(entities.Expense), args.Error(1)
}

func (m *ExpensesRepositoryMock) Create(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	args := m.Called(ctx, expense)

	return args.Get(0).(entities.Expense), args.Error(1)
}

func (m *ExpensesRepositoryMock) Update(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	args := m.Called(ctx, expense)

	return args.Get(0).(entities.Expense), args.Error(1)
}

func (m *ExpensesRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func NewExpensesRepositoryMock() *ExpensesRepositoryMock {
	return &ExpensesRepositoryMock{}
}
