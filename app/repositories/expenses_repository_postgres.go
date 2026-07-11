package repositories

import (
	"context"
	"database/sql"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/json"
	"github.com/ZioCastoro/restaurant_assignment/postgre"
)

const createExpenseQuery = `SELECT create_expense(:expense);`
const updateExpenseQuery = `SELECT update_expense(:expense);`
const deleteExpenseQuery = `SELECT delete_expense(:expense_id);`
const listExpensesQuery = `SELECT list_expenses(:query_text, :query_supplier, :query_from, :query_to, :query_limit, :query_offset);`
const countExpensesQuery = `SELECT count_expenses(:query_text, :query_supplier, :query_from, :query_to);`
const findExpenseQuery = `SELECT find_expense(:expense_id);`

type ExpensesRepositoryPostgres struct {
	db *sql.DB
}

func (e *ExpensesRepositoryPostgres) List(
	ctx context.Context,
	filters entities.ExpenseFilters,
) (entities.PaginatedExpenses, error) {
	pqResults, err := postgre.PrepareQuery(
		listExpensesQuery,
		postgre.Named("query_text", filters.Text),
		postgre.Named("query_supplier", filters.Supplier),
		postgre.Named("query_from", filters.From),
		postgre.Named("query_to", filters.To),
		postgre.Named("query_limit", filters.GetLimit()),
		postgre.Named("query_offset", filters.GetOffset()),
	)
	if err != nil {
		return entities.PaginatedExpenses{}, err
	}

	pqTotal, err := postgre.PrepareQuery(
		countExpensesQuery,
		postgre.Named("query_text", filters.Text),
		postgre.Named("query_supplier", filters.Supplier),
		postgre.Named("query_from", filters.From),
		postgre.Named("query_to", filters.To),
	)
	if err != nil {
		return entities.PaginatedExpenses{}, err
	}

	results, err := postgre.QueryJSON[entities.Expense](ctx, e.db, pqResults)
	if err != nil {
		return entities.PaginatedExpenses{}, err
	}

	total, err := postgre.QueryRow[int64](ctx, e.db, pqTotal)
	if err != nil {
		return entities.PaginatedExpenses{}, err
	}

	return entities.PaginatedExpenses{
		Results: results,
		Total:   total,
	}, nil
}

func (e *ExpensesRepositoryPostgres) Find(
	ctx context.Context,
	id string,
) (entities.Expense, error) {
	pq, err := postgre.PrepareQuery(
		findExpenseQuery,
		postgre.Named("expense_id", id),
	)
	if err != nil {
		return entities.Expense{}, err
	}

	return postgre.QueryRowJSON[entities.Expense](ctx, e.db, pq)
}

func (e *ExpensesRepositoryPostgres) Create(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	pq, err := postgre.PrepareQuery(
		createExpenseQuery,
		postgre.Named("expense", json.MustMarshal(expense)),
	)
	if err != nil {
		return entities.Expense{}, err
	}

	return postgre.QueryRowJSON[entities.Expense](ctx, e.db, pq)
}

func (e *ExpensesRepositoryPostgres) Update(
	ctx context.Context,
	expense entities.Expense,
) (entities.Expense, error) {
	pq, err := postgre.PrepareQuery(
		updateExpenseQuery,
		postgre.Named("expense", json.MustMarshal(expense)),
	)
	if err != nil {
		return entities.Expense{}, err
	}

	return postgre.QueryRowJSON[entities.Expense](ctx, e.db, pq)
}

func (e *ExpensesRepositoryPostgres) Delete(
	ctx context.Context,
	id string,
) error {
	pq, err := postgre.PrepareQuery(
		deleteExpenseQuery,
		postgre.Named("expense_id", id),
	)
	if err != nil {
		return err
	}

	_, err = postgre.Exec(ctx, e.db, pq)
	if err != nil {
		return err
	}

	return nil
}

var _ ExpensesRepositoryInterface = new(ExpensesRepositoryPostgres)

func NewExpensesRepositoryPostgres(db *sql.DB) *ExpensesRepositoryPostgres {
	return &ExpensesRepositoryPostgres{
		db: db,
	}
}
