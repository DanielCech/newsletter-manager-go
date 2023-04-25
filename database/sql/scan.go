package sql

import (
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
)

const (
	singleEntity = 1
)

var ErrExecNotOne = errors.New("not found exactly 1 row for exec")

// Read is a helper function which calls a database query and expects a single result.
// Returned value is a pointer type.
// Returns an error matchable by IsNotFound if no rows were returned.
func Read[T any](dctx DataContext, query string, args ...any) (*T, error) {
	var result T
	if err := pgxscan.Get(dctx.Ctx(), querier(dctx), &result, query, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

// ReadValue is a helper function which calls a database query and expects a single result.
// Returns an error matchable by IsNotFound if no rows were returned.
func ReadValue[T any](dctx DataContext, query string, args ...any) (T, error) {
	var result T
	if err := pgxscan.Get(dctx.Ctx(), querier(dctx), &result, query, args...); err != nil {
		return result, err
	}
	return result, nil
}

// List is a helper function which calls a database query and expects a list of results.
func List[T any](dctx DataContext, query string, args ...any) ([]T, error) {
	var result []T
	if err := pgxscan.Select(dctx.Ctx(), querier(dctx), &result, query, args...); err != nil {
		return nil, err
	}
	return result, nil
}

// ExecOne is a helper function which calls a database exec.
// Expected is the single affected row.
// Returns an error matchable by IsNotFound if Exec did not affect 1 row.
func ExecOne(dctx DataContext, query string, args ...any) error {
	r, err := querier(dctx).Exec(dctx.Ctx(), query, args...)
	if err != nil {
		return err
	}
	if affected := r.RowsAffected(); affected != singleEntity {
		return fmt.Errorf("%w: expected 1 row affected but affected %d rows", ErrExecNotOne, affected)
	}
	return nil
}

// Exec is a helper function which calls a database exec.
func Exec(dctx DataContext, query string, args ...any) error {
	if _, err := querier(dctx).Exec(dctx.Ctx(), query, args...); err != nil {
		return err
	}
	return nil
}
