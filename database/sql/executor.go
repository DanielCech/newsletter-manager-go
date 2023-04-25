package sql

import (
	"context"
	"errors"
	"fmt"
)

// DataSource defines common functions for manipulating with database specific features.
type DataSource interface {
	AcquireConnCtx(ctx context.Context) (DataContext, error)
	ReleaseConnCtx(dctx DataContext) error
	Begin(ctx context.Context) (DataContext, error)
	Commit(dctx DataContext) error
	Rollback(dctx DataContext) error
}

// WithConnection is a helper function which acquires and releases connection for your workload function.
func WithConnection(ctx context.Context, s DataSource, f func(DataContext) error) (err error) {
	dctx, err := s.AcquireConnCtx(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer func() {
		if rErr := s.ReleaseConnCtx(dctx); rErr != nil {
			err = errors.Join(rErr, err)
		}
	}()
	return f(dctx)
}

// WithConnectionResult is a helper function which acquires and releases connection for your workload function.
func WithConnectionResult[T any](ctx context.Context, s DataSource, f func(DataContext) (T, error)) (result T, err error) {
	dctx, err := s.AcquireConnCtx(ctx)
	if err != nil {
		return result, fmt.Errorf("acquiring connection: %w", err)
	}
	defer func() {
		if rErr := s.ReleaseConnCtx(dctx); rErr != nil {
			err = errors.Join(rErr, err)
		}
	}()
	return f(dctx)
}

// WithTransaction is a helper function which begins/commits/rollbacks transaction for your workload function.
func WithTransaction(ctx context.Context, s DataSource, f func(DataContext) error) (err error) {
	dctx, err := s.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if rErr := s.Rollback(dctx); rErr != nil {
			err = errors.Join(rErr, err)
		}
	}()

	if err = f(dctx); err != nil {
		return err
	}

	if err = s.Commit(dctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// WithTransactionResult is a helper function which begins/commits/rollbacks transaction for your workload function.
func WithTransactionResult[T any](ctx context.Context, s DataSource, f func(DataContext) (T, error)) (result T, err error) {
	dctx, err := s.Begin(ctx)
	if err != nil {
		return result, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if rErr := s.Rollback(dctx); rErr != nil {
			err = errors.Join(rErr, err)
		}
	}()

	result, err = f(dctx)
	if err != nil {
		return result, err
	}

	if err = s.Commit(dctx); err != nil {
		return result, fmt.Errorf("committing transaction: %w", err)
	}

	return result, nil
}
