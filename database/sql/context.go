package sql

import (
	"context"
)

var (
	contextKey = struct {
		querier ctxKeyQuerier
	}{}
)

type (
	ctxKeyQuerier struct{}
)

// DataContext handles internally low level database logic.
// Such context can contain for example a connection or transaction.
type DataContext interface {
	// Ctx returns the underlying context containing internal components of used database type.
	Ctx() context.Context
}

// WithQuerier passes querier to the context.
func WithQuerier(ctx context.Context, querier Querier) context.Context {
	return context.WithValue(ctx, contextKey.querier, querier)
}

// QuerierFromCtx gets DataContext from the context.
func QuerierFromCtx(ctx context.Context) (Querier, bool) {
	querier, ok := ctx.Value(contextKey.querier).(Querier)
	return querier, ok
}

// dbContext wraps context.Context to distinguish between regular and database context.
type dbContext struct {
	context.Context
}

// Ctx returns the underlying context containing pgx connection or transaction.
func (d dbContext) Ctx() context.Context {
	return d.Context
}
