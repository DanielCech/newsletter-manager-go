package context

import (
	"context"

	"newsletter-manager-go/types/id"
)

var (
	contextKey = struct {
		AuthorID ctxKeyAuthorID
		userRole ctxKeyUserRole
	}{}
)

type (
	ctxKeyAuthorID struct{}
	ctxKeyUserRole struct{}
)

// WithAuthorID passes user ID to the context.
func WithAuthorID(ctx context.Context, authorID id.Author) context.Context {
	return context.WithValue(ctx, contextKey.AuthorID, authorID)
}

// AuthorIDFromCtx gets user ID from the context.
func AuthorIDFromCtx(ctx context.Context) (id.Author, bool) {
	authorID, ok := ctx.Value(contextKey.AuthorID).(id.Author)
	return authorID, ok
}
