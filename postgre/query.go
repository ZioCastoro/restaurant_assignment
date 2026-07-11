package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ZioCastoro/restaurant_assignment/json"
)

var (
	ErrMissingParam   = errors.New("missing param(s) for query")
	ErrDuplicateParam = errors.New("duplicate param for query")
	ErrParamNotFound  = errors.New("param name not found in query")
	ErrNULLJSON       = errors.New("query returned NULL, expected JSON")
)

type PreparedQuery struct {
	Query string
	Args  []any
}

// QueryRow executes a query that is expected to return at most one row.
// Wraps the QueryRowContext method of the Driver interface.
// The result is scanned into a value of type T and returned.
// T must be a type that can be scanned by sql.Row.Scan.
// If no rows are returned, sql.ErrNoRows is returned.
// This method is intended to be used for queries that must return exactly one row.
func QueryRow[T any](
	ctx context.Context,
	d Driver,
	pq PreparedQuery,
) (T, error) {
	var res T

	row := d.QueryRowContext(ctx, pq.Query, pq.Args...)

	if err := row.Err(); err != nil {
		return *new(T), fmt.Errorf("row.Err(): %w", err)
	}

	if err := row.Scan(&res); err != nil {
		return *new(T), fmt.Errorf("row.Scan(): %w", err)
	}

	return res, nil
}

// Query executes a query that returns multiple rows.
// Wraps the QueryContext method of the Driver interface.
// The results are scanned into a slice of values of type T and returned.
// T must be a type that can be scanned by sql.Rows.Scan.
// If no rows are returned, an empty slice is returned and no error.
// This method is intended to be used for queries that can return zero or more rows.
func Query[T any](
	ctx context.Context,
	d Driver,
	pq PreparedQuery,
) ([]T, error) {
	rows, err := d.QueryContext(ctx, pq.Query, pq.Args...)
	if err != nil {
		return nil, fmt.Errorf("QueryContext(): %w", err)
	}
	defer rows.Close()

	var res []T

	for rows.Next() {
		var item T
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		res = append(res, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return res, nil
}

// QueryRowJSON executes a query that is expected to return at most one row with a JSON result.
// Wraps the QueryRowContext method of the Driver interface.
// The result is unmarshaled into a value of type T and returned.
// If no rows are returned, sql.ErrNoRows is returned.
// If the returned JSON is NULL, ErrNullJSON is returned.
// This method is intended to be used for queries that must return exactly one row with JSON data.
func QueryRowJSON[T any](
	ctx context.Context,
	d Driver,
	pq PreparedQuery,
) (T, error) {
	// use sql.NullString because raw bytes are not supported by Scan.
	resJSON, err := QueryRow[sql.NullString](ctx, d, pq)
	if err != nil {
		return *new(T), err
	}

	if !resJSON.Valid {
		return *new(T), ErrNULLJSON
	}

	res, err := json.Unmarshal[T]([]byte(resJSON.String))
	if err != nil {
		return *new(T), fmt.Errorf("json.Unmarshal(): %w", err)
	}

	return res, nil
}

// QueryJSON executes a query that returns multiple rows with JSON results.
// Wraps the QueryContext method of the Driver interface.
// The results are unmarshaled into a slice of values of type T and returned.
// If no rows are returned, an empty slice is returned and no error.
// If any of the returned JSON values is NULL, ErrNullJSON is returned.
// This method is intended to be used for queries that can return zero or more rows with JSON data.
func QueryJSON[T any](
	ctx context.Context,
	d Driver,
	pq PreparedQuery,
) ([]T, error) {
	// use sql.NullString because raw bytes are not supported by Scan.
	resJSON, err := Query[sql.NullString](ctx, d, pq)
	if err != nil {
		return nil, err
	}

	if len(resJSON) == 0 {
		return []T(nil), nil
	}

	res := make([]T, 0, len(resJSON))

	for _, itemJSON := range resJSON {
		if !itemJSON.Valid {
			return nil, ErrNULLJSON
		}

		item, err := json.Unmarshal[T]([]byte(itemJSON.String))
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal(): %w", err)
		}

		res = append(res, item)
	}

	return res, nil
}

// Exec executes a query that doesn't return rows.
// Wraps the ExecContext method of the Driver interface.
func Exec(
	ctx context.Context,
	d Driver,
	pq PreparedQuery,
) (sql.Result, error) {
	res, err := d.ExecContext(ctx, pq.Query, pq.Args...)
	if err != nil {
		return nil, fmt.Errorf("ExecContext(): %w", err)
	}

	return res, nil
}

// PrepareQuery converts a query with named params having ':param' as placeholder
// into a query with sequential params ($1, $2, ...).
// It also converts QueryArgs into a slice of arguments.
// Returns ErrDuplicateParam if a param is defined more than once.
// Returns ErrParamNotFound if a param name is not found in the query.
// Returns ErrMissingParam if not all params were replaced.
func PrepareQuery(query string, namedArgs ...NamedArg) (PreparedQuery, error) {
	argsMap := make(map[string]any)
	for _, arg := range namedArgs {
		if _, exists := argsMap[arg.Name()]; exists {
			return PreparedQuery{}, fmt.Errorf("%w: %s", ErrDuplicateParam, arg.Name())
		}

		argsMap[arg.Name()] = arg.Value()
	}

	argSlice := make([]any, 0, len(namedArgs))

	i := 0

	for param, value := range argsMap {
		i++
		var err error

		query, err = replaceParam(query, param, fmt.Sprintf("$%d", i))
		if err != nil {
			return PreparedQuery{}, err
		}

		argSlice = append(argSlice, value)
	}

	// check if all named params were replaced.
	if strings.Contains(query, ":") {
		return PreparedQuery{}, fmt.Errorf("%w: %s", ErrMissingParam, query)
	}

	return PreparedQuery{
		Query: query,
		Args:  argSlice,
	}, nil
}

// replaceParam replaces a param in a query with a given value.
// Returns ErrParamNotFound if the param name was not found in the query.
func replaceParam(query, param, value string) (string, error) {
	indexes := regexp.MustCompile(":"+param+`(?:([,;\s)]))`).FindAllStringIndex(query, -1)
	if indexes == nil {
		return "", fmt.Errorf("%w: %s", ErrParamNotFound, param)
	}

	indexes = append(indexes, []int{len(query)})

	var res strings.Builder
	res.WriteString(query[:indexes[0][0]])

	for j := 1; j < len(indexes); j++ {
		res.WriteString(value + query[indexes[j-1][1]-1:indexes[j][0]])
	}

	return res.String(), nil
}
