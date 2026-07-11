package postgre_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ZioCastoro/restaurant_assignment/postgre"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDB struct {
	*sql.DB
	mock sqlmock.Sqlmock
}

type result struct {
	Field int `json:"field"`
}

func Test_QueryRow(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		testDB testDB
		pq     postgre.PreparedQuery
	}

	// define test cases.
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryRow(t, 1),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "query error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryError(t, "some error"),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "some error",
		},
		{
			name: "scan error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryRow(t, "result"),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "sql: Scan error on column",
		},
		{
			name: "no rows error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryNoRows(t),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "sql: no rows in result set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				if tt.args.testDB.DB != nil {
					_ = tt.args.testDB.Close()
				}
			})

			var err error
			var got int
			// call method under test with panic recover.
			got, err = func(args args) (int, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.QueryRow[int](args.ctx, args.testDB, args.pq)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Equal(t, tt.want, got)

			// assert mock expectations.
			assert.NoError(t, tt.args.testDB.mock.ExpectationsWereMet())
		})
	}
}

func Test_Query(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		testDB testDB
		pq     postgre.PreparedQuery
	}

	// define test cases.
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{1, 2, 3}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name: "happy flow empty result",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryNoRows(t),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    []int(nil),
			wantErr: false,
		},
		{
			name: "query error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryError(t, "some error"),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []int(nil),
			wantErr: true,
			errMsg:  "some error",
		},
		{
			name: "scan error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{"result"}),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    []int(nil),
			wantErr: true,
			errMsg:  "sql: Scan error on column",
		},
		{
			name: "rows close error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryCloseError(t, []any{1}, "close error"),
				pq: postgre.PreparedQuery{
					Query: "any query",
					Args:  nil,
				},
			},
			want:    []int(nil),
			wantErr: true,
			errMsg:  "close error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				if tt.args.testDB.DB != nil {
					_ = tt.args.testDB.Close()
				}
			})

			var err error
			var got []int
			// call method under test with panic recover.
			got, err = func(args args) ([]int, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.Query[int](args.ctx, args.testDB, args.pq)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Equal(t, tt.want, got)

			// assert mock expectations.
			assert.NoError(t, tt.args.testDB.mock.ExpectationsWereMet())
		})
	}
}

func Test_QueryRowJSON(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		testDB testDB
		pq     postgre.PreparedQuery
	}

	// define test cases.
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{`{"field":1}`}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    result{Field: 1},
			wantErr: false,
		},
		{
			name: "no rows error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryNoRows(t),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    result{},
			wantErr: true,
			errMsg:  "sql: no rows in result set",
		},
		{
			name: "NULL JSON result error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{nil}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    result{},
			wantErr: true,
			errMsg:  "query returned NULL, expected JSON",
		},
		{
			name: "invalid JSON error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{`invalid json`}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    result{},
			wantErr: true,
			errMsg:  "json.Unmarshal():",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				if tt.args.testDB.DB != nil {
					_ = tt.args.testDB.Close()
				}
			})

			var err error
			var got result
			// call method under test with panic recover.
			got, err = func(args args) (result, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.QueryRowJSON[result](args.ctx, args.testDB, args.pq)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Equal(t, tt.want, got)

			// assert mock expectations.
			assert.NoError(t, tt.args.testDB.mock.ExpectationsWereMet())
		})
	}
}

func Test_QueryJSON(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		testDB testDB
		pq     postgre.PreparedQuery
	}

	// define test cases.
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{`{"field":1}`, `{"field":2}`}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []result{{Field: 1}, {Field: 2}},
			wantErr: false,
		},
		{
			name: "happy flow empty result",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryNoRows(t),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []result(nil),
			wantErr: false,
		},
		{
			name: "query error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQueryError(t, "some error"),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []result(nil),
			wantErr: true,
			errMsg:  "some error",
		},
		{
			name: "NULL JSON result error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{nil}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []result(nil),
			wantErr: true,
			errMsg:  "query returned NULL, expected JSON",
		},
		{
			name: "invalid JSON error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockQuery(t, []any{`invalid json`}),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			want:    []result(nil),
			wantErr: true,
			errMsg:  "json.Unmarshal():",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				if tt.args.testDB.DB != nil {
					_ = tt.args.testDB.Close()
				}
			})

			var err error
			var got []result
			// call method under test with panic recover.
			got, err = func(args args) ([]result, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.QueryJSON[result](args.ctx, args.testDB, args.pq)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Equal(t, tt.want, got)

			// assert mock expectations.
			assert.NoError(t, tt.args.testDB.mock.ExpectationsWereMet())
		})
	}
}

func Test_Exec(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		testDB testDB
		pq     postgre.PreparedQuery
	}

	// define test cases.
	tests := []struct {
		name             string
		args             args
		wantLastInsertID int64
		wantRowsAffected int64
		wantErr          bool
		errMsg           string
	}{
		{
			name: "happy flow",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockExec(t, 1, 1),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			wantLastInsertID: 1,
			wantRowsAffected: 1,
			wantErr:          false,
		},
		{
			name: "exec error",
			args: args{
				ctx:    context.Background(),
				testDB: setDriverMockExecFail(t, "some exec error"),
				pq: postgre.PreparedQuery{
					Query: "any query",
				},
			},
			wantLastInsertID: 1,
			wantRowsAffected: 1,
			wantErr:          true,
			errMsg:           "some exec error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				if tt.args.testDB.DB != nil {
					_ = tt.args.testDB.Close()
				}
			})

			var err error
			var got sql.Result
			// call method under test with panic recover.
			got, err = func(args args) (sql.Result, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.Exec(args.ctx, args.testDB, args.pq)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			// we need to split the logic between error and non-error cases
			// because we can't compare the full sql.Result interface.
			// and we need to handle the nil case differently to avoid panic issues.
			if !tt.wantErr {
				require.NotNil(t, got)

				lastInsertID, err := got.LastInsertId()
				require.NoError(t, err)
				assert.Equal(t, tt.wantLastInsertID, lastInsertID)

				rowsAffected, err := got.RowsAffected()
				require.NoError(t, err)
				assert.Equal(t, tt.wantRowsAffected, rowsAffected)
			} else {
				require.Nil(t, got)
			}

			// assert mock expectations.
			assert.NoError(t, tt.args.testDB.mock.ExpectationsWereMet())
		})
	}
}

func Test_PrepareQuery(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		query     string
		namedArgs []postgre.NamedArg
	}

	// define test cases.
	tests := []struct {
		name    string
		args    args
		want    []postgre.PreparedQuery
		wantErr bool
		errMsg  string
	}{
		{
			name: "happy flow simple query",
			args: args{
				query: "SELECT * FROM users WHERE id = :id;",
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
				},
			},
			want: []postgre.PreparedQuery{
				{
					Query: "SELECT * FROM users WHERE id = $1;",
					Args:  []any{42},
				},
			},
			wantErr: false,
		},
		{
			name: "happy flow function call",
			args: args{
				query: "SELECT find_user(:id, :name);",
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
					postgre.Named("name", "Alice"),
				},
			},
			want: []postgre.PreparedQuery{
				{
					Query: "SELECT find_user($1, $2);",
					Args:  []any{42, "Alice"},
				},
				{
					Query: "SELECT find_user($2, $1);",
					Args:  []any{"Alice", 42},
				},
			},
			wantErr: false,
		},
		{
			name: "happy flow multiple columns",
			args: args{
				query: `SELECT find_user(
					:id,
					:name
				);`,
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
					postgre.Named("name", "Alice"),
				},
			},
			want: []postgre.PreparedQuery{
				{
					Query: `SELECT find_user(
					$1,
					$2
				);`,
					Args: []any{42, "Alice"},
				},
				{
					Query: `SELECT find_user(
					$2,
					$1
				);`,
					Args: []any{"Alice", 42},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate parameter",
			args: args{
				query: "SELECT * FROM users WHERE id = :id;",
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
					postgre.Named("id", 43),
				},
			},
			want:    []postgre.PreparedQuery{{}},
			wantErr: true,
			errMsg:  "duplicate param for query",
		},
		{
			name: "parameter not found",
			args: args{
				query: "SELECT * FROM users WHERE name = :name;",
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
				},
			},
			want:    []postgre.PreparedQuery{{}},
			wantErr: true,
			errMsg:  "param name not found in query",
		},
		{
			name: "missing parameter",
			args: args{
				query: "SELECT find_user(:id, :name);",
				namedArgs: []postgre.NamedArg{
					postgre.Named("id", 42),
				},
			},
			want:    []postgre.PreparedQuery{{}},
			wantErr: true,
			errMsg:  "missing param(s) for query",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var got postgre.PreparedQuery
			// call method under test with panic recover.
			got, err = func(args args) (postgre.PreparedQuery, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.PrepareQuery(args.query, args.namedArgs...)
			}(tt.args)

			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			assert.Contains(t, tt.want, got, "prepared query not in expected list")
		})
	}
}

func setDriverMockQueryRow(
	t *testing.T,
	response any,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"res"}).AddRow(response))

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockQueryError(
	t *testing.T,
	msg string,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectQuery("").WillReturnError(errors.New(msg))

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockQuery(
	t *testing.T,
	responses []any,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	rows := sqlmock.NewRows([]string{"res"})
	for _, response := range responses {
		rows.AddRow(response)
	}

	mock.ExpectQuery("").WillReturnRows(rows)

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockQueryCloseError(
	t *testing.T,
	responses []any,
	msg string,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	rows := sqlmock.NewRows([]string{"res"})
	for _, response := range responses {
		rows.AddRow(response)
	}

	rows.CloseError(errors.New(msg))

	mock.ExpectQuery("").WillReturnRows(rows)

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockQueryNoRows(t *testing.T) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"res"}))

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockExec(
	t *testing.T,
	lastInsertID int64,
	rowsAffected int64,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(lastInsertID, rowsAffected))

	return testDB{
		DB:   db,
		mock: mock,
	}
}

func setDriverMockExecFail(
	t *testing.T,
	msg string,
) testDB {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectExec("").WillReturnError(errors.New(msg))

	return testDB{
		DB:   db,
		mock: mock,
	}
}
