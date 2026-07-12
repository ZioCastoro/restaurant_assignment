package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/app/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NewExpensesService(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// define test cases.
	tests := []struct {
		name   string
		args   args
		errMsg string
	}{
		{
			name: "happy flow",
			args: args{
				repository: repositories.NewExpensesRepositoryMock(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var got *ExpensesService
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				got = NewExpensesService(args.repository)
			}(tt.args)

			// assert no panic.
			require.NoError(t, err)

			// assert result.
			require.NotNil(t, got)
			assert.NotNil(t, got.repository)

			// assert mock expectations.
			tt.args.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_List(t *testing.T) {
	t.Parallel()

	updatedAt := time.Now()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx     context.Context
		filters entities.ExpenseFilters
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.PaginatedExpenses
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockList(
					t,
					getTestExpenses(t, updatedAt),
					nil,
				),
			},
			args: args{
				ctx:     context.Background(),
				filters: entities.ExpenseFilters{},
			},
			want:    getTestExpenses(t, updatedAt),
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockList(
					t,
					entities.PaginatedExpenses{},
					errors.New("list error"),
				),
			},
			args: args{
				ctx:     context.Background(),
				filters: entities.ExpenseFilters{},
			},
			wantErr: true,
			errMsg:  "list error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx:     canceledCtx,
				filters: entities.ExpenseFilters{},
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			var got entities.PaginatedExpenses
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				got, err = service.List(args.ctx, args.filters)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				assert.Zero(t, got)
			} else {
				require.NoError(t, err)

				// assert result.
				assert.Equal(t, tt.want, got)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_Find(t *testing.T) {
	t.Parallel()

	updatedAt := time.Now()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx context.Context
		id  string
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Expense
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockFind(
					t,
					"exp-1",
					getTestExpense(t, updatedAt),
					nil,
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "exp-1",
			},
			want:    getTestExpense(t, updatedAt),
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockFind(
					t,
					"exp-1",
					entities.Expense{},
					errors.New("find error"),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "exp-1",
			},
			wantErr: true,
			errMsg:  "find error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx: canceledCtx,
				id:  "exp-1",
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			var got entities.Expense
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				got, err = service.Find(args.ctx, args.id)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				assert.Zero(t, got)
			} else {
				require.NoError(t, err)

				// assert result.
				assert.Equal(t, tt.want, got)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_Create(t *testing.T) {
	t.Parallel()

	updatedAt := time.Now()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx     context.Context
		expense entities.Expense
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Expense
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockCreate(
					t,
					getTestExpense(t, updatedAt),
					nil,
				),
			},
			args: args{
				ctx:     context.Background(),
				expense: getTestExpenseInput(t, updatedAt),
			},
			want:    getTestExpense(t, updatedAt),
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockCreate(
					t,
					entities.Expense{},
					errors.New("create error"),
				),
			},
			args: args{
				ctx:     context.Background(),
				expense: getTestExpenseInput(t, updatedAt),
			},
			wantErr: true,
			errMsg:  "create error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx:     canceledCtx,
				expense: getTestExpenseInput(t, updatedAt),
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			var got entities.Expense
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				got, err = service.Create(args.ctx, args.expense)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				assert.Zero(t, got)
			} else {
				require.NoError(t, err)

				// assert result.
				assert.Equal(t, tt.want, got)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_Update(t *testing.T) {
	t.Parallel()

	updatedAt := time.Now()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx     context.Context
		expense entities.Expense
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Expense
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockUpdate(
					t,
					getTestExpense(t, updatedAt),
					nil,
				),
			},
			args: args{
				ctx:     context.Background(),
				expense: getTestExpense(t, updatedAt),
			},
			want:    getTestExpense(t, updatedAt),
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockUpdate(
					t,
					entities.Expense{},
					errors.New("update error"),
				),
			},
			args: args{
				ctx:     context.Background(),
				expense: getTestExpense(t, updatedAt),
			},
			wantErr: true,
			errMsg:  "update error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx:     canceledCtx,
				expense: getTestExpense(t, updatedAt),
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			var got entities.Expense
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				got, err = service.Update(args.ctx, args.expense)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				assert.Zero(t, got)
			} else {
				require.NoError(t, err)

				// assert result.
				assert.Equal(t, tt.want, got)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_Delete(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx context.Context
		id  string
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockDelete(t, "exp-1", nil),
			},
			args: args{
				ctx: context.Background(),
				id:  "exp-1",
			},
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockDelete(t, "exp-1", errors.New("delete error")),
			},
			args: args{
				ctx: context.Background(),
				id:  "exp-1",
			},
			wantErr: true,
			errMsg:  "delete error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx: canceledCtx,
				id:  "exp-1",
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				err = service.Delete(args.ctx, args.id)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func Test_ExpensesService_Import(t *testing.T) {
	t.Parallel()

	updatedAt := time.Now()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// initialize service fields.
	type fields struct {
		repository *repositories.ExpensesRepositoryMock
	}

	// initialize input arguments.
	type args struct {
		ctx      context.Context
		expenses []entities.Expense
	}

	// define test cases.
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			fields: fields{
				repository: setRepositoryMockImport(
					t,
					getTestExpense(t, updatedAt),
					nil,
					2,
				),
			},
			args: args{
				ctx:      context.Background(),
				expenses: getTestExpenseInputs(t, updatedAt),
			},
			wantErr: false,
		},
		{
			name: "repository error",
			fields: fields{
				repository: setRepositoryMockImport(
					t,
					entities.Expense{},
					errors.New("create error"),
					2,
				),
			},
			args: args{
				ctx:      context.Background(),
				expenses: getTestExpenseInputs(t, updatedAt),
			},
			wantErr: true,
			errMsg:  "create error",
		},
		{
			name: "cancelled context",
			fields: fields{
				repository: repositories.NewExpensesRepositoryMock(),
			},
			args: args{
				ctx:      canceledCtx,
				expenses: getTestExpenseInputs(t, updatedAt),
			},
			wantErr: true,
			errMsg:  context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &ExpensesService{
				repository: tt.fields.repository,
			}

			var err error
			// call method under test with panic recover.
			func(args args) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				err = service.Import(args.ctx, args.expenses)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert mock expectations.
			tt.fields.repository.AssertExpectations(t)
		})
	}
}

func getTestExpense(
	t *testing.T,
	updatedAt time.Time,
) entities.Expense {
	t.Helper()

	return entities.Expense{
		ID:          "exp-1",
		Supplier:    "Acme",
		IssuedAt:    updatedAt,
		Description: "Office supplies",
		Qty:         2,
		CostPerUnit: 10.50,
		Category:    "supplies",
	}
}

func getTestExpenseInput(
	t *testing.T,
	updatedAt time.Time,
) entities.Expense {
	t.Helper()

	expense := getTestExpense(t, updatedAt)
	expense.ID = ""

	return expense
}

func getTestExpenseInputs(
	t *testing.T,
	updatedAt time.Time,
) []entities.Expense {
	t.Helper()

	secondExpense := getTestExpense(t, updatedAt)
	secondExpense.ID = ""
	secondExpense.Supplier = "Globex"
	secondExpense.Description = "Kitchen equipment"
	secondExpense.Qty = 1
	secondExpense.CostPerUnit = 250.00
	secondExpense.Category = "equipment"

	return []entities.Expense{
		getTestExpenseInput(t, updatedAt),
		secondExpense,
	}
}

func getTestExpenses(
	t *testing.T,
	updatedAt time.Time,
) entities.PaginatedExpenses {
	t.Helper()

	return entities.PaginatedExpenses{
		Results: []entities.Expense{
			getTestExpense(t, updatedAt),
			{
				ID:          "exp-2",
				Supplier:    "Globex",
				IssuedAt:    updatedAt,
				Description: "Kitchen equipment",
				Qty:         1,
				CostPerUnit: 250.00,
				Category:    "equipment",
			},
		},
		Total: 2,
	}
}

func setRepositoryMockList(
	t *testing.T,
	res entities.PaginatedExpenses,
	err error,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	mockRepository.On(
		"List",
		mock.Anything,
		mock.AnythingOfType("entities.ExpenseFilters"),
	).Return(res, err)

	return mockRepository
}

func setRepositoryMockFind(
	t *testing.T,
	id string,
	res entities.Expense,
	err error,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	mockRepository.On(
		"Find",
		mock.Anything,
		id,
	).Return(res, err)

	return mockRepository
}

func setRepositoryMockCreate(
	t *testing.T,
	res entities.Expense,
	err error,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	mockRepository.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("entities.Expense"),
	).Return(res, err)

	return mockRepository
}

func setRepositoryMockUpdate(
	t *testing.T,
	res entities.Expense,
	err error,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	mockRepository.On(
		"Update",
		mock.Anything,
		mock.AnythingOfType("entities.Expense"),
	).Return(res, err)

	return mockRepository
}

func setRepositoryMockDelete(
	t *testing.T,
	id string,
	err error,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	mockRepository.On(
		"Delete",
		mock.Anything,
		id,
	).Return(err)

	return mockRepository
}

func setRepositoryMockImport(
	t *testing.T,
	res entities.Expense,
	err error,
	times int,
) *repositories.ExpensesRepositoryMock {
	t.Helper()

	mockRepository := repositories.NewExpensesRepositoryMock()
	if times == 0 {
		return mockRepository
	}

	mockRepository.On(
		"Create",
		mock.Anything,
		mock.AnythingOfType("entities.Expense"),
	).Return(res, err).Times(times)

	return mockRepository
}
