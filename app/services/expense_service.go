package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/app/repositories"
	"github.com/gofiber/fiber/v2/log"
)

type ExpensesService struct {
	repository repositories.ExpensesRepositoryInterface
}

func (e *ExpensesService) List(
	ctx context.Context,
	filters entities.ExpenseFilters,
) (entities.PaginatedExpenses, error) {
	if ctx.Err() != nil {
		return entities.PaginatedExpenses{}, ctx.Err()
	}

	res, err := e.repository.List(ctx, filters)
	if err != nil {
		log.Errorf("resource: expenses, action: list, error: %s", err.Error())

		return entities.PaginatedExpenses{}, err
	}

	return res, nil
}

func (e *ExpensesService) Find(
	ctx context.Context,
	id string,
) (entities.Expense, error) {
	if ctx.Err() != nil {
		return entities.Expense{}, ctx.Err()
	}

	res, err := e.repository.Find(ctx, id)
	if err != nil {
		log.Errorf("resource: expenses, action: find, id: %s, error: %s", id, err.Error())

		return entities.Expense{}, err
	}

	return res, nil
}

func (e *ExpensesService) Create(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	if ctx.Err() != nil {
		return entities.Expense{}, ctx.Err()
	}

	expense.GenerateID()

	res, err := e.repository.Create(ctx, expense)
	if err != nil {
		log.Errorf("resource: expenses, action: create, error: %s", err.Error())

		return entities.Expense{}, err
	}

	return res, nil
}

func (e *ExpensesService) Update(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	if ctx.Err() != nil {
		return entities.Expense{}, ctx.Err()
	}

	res, err := e.repository.Update(ctx, expense)
	if err != nil {
		log.Errorf("resource: expenses, action: update, id: %s, error: %s", expense.ID, err.Error())

		return entities.Expense{}, err
	}

	return res, nil
}

func (e *ExpensesService) Delete(
	ctx context.Context,
	id string,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := e.repository.Delete(ctx, id); err != nil {
		log.Errorf("resource: expenses, action: delete, id: %s, error: %s", id, err.Error())

		return err
	}

	return nil
}

func (e *ExpensesService) Import(
	ctx context.Context,
	expenses []entities.Expense,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, expense := range expenses {
		wg.Add(1)

		go func(expense entities.Expense) {
			defer wg.Done()

			if _, err := e.Create(ctx, expense); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(expense)
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		log.Errorf("resource: expenses, action: import, error: %s", err.Error())

		return err
	}
	// TODO: send a notification to the user with the result of the import.

	return nil
}

var _ ExpenseServiceInterface = new(ExpensesService)

func NewExpensesService(
	repository repositories.ExpensesRepositoryInterface,
) *ExpensesService {
	return &ExpensesService{
		repository: repository,
	}
}
