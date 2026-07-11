package postgre

import (
	"context"
	"database/sql"
)

var (
	_ Driver = (*sql.DB)(nil)
	_ Driver = (*sql.Tx)(nil)
)

type Driver interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
