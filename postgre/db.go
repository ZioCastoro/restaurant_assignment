package postgre

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver (needed for sql.Open)
)

func Connect(ctx context.Context, driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(): %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext(): %w", err)
	}

	return db, nil
}
