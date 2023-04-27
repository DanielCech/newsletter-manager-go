package context

import (
	"context"

	domauthor "newsletter-manager-go/domain/author"
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
func WithAuthorID(ctx context.Context, AuthorID id.Author) context.Context {
	return context.WithValue(ctx, contextKey.AuthorID, AuthorID)
}

// AuthorIDFromCtx gets user ID from the context.
func AuthorIDFromCtx(ctx context.Context) (id.Author, bool) {
	AuthorID, ok := ctx.Value(contextKey.AuthorID).(id.Author)
	return AuthorID, ok
}

// WithUserRole passes user role to the context.
func WithUserRole(ctx context.Context, role domauthor.Role) context.Context {
	return context.WithValue(ctx, contextKey.userRole, role)
}

// UserRoleFromCtx gets user role from the context.
func UserRoleFromCtx(ctx context.Context) (domauthor.Role, bool) {
	userRole, ok := ctx.Value(contextKey.userRole).(domauthor.Role)
	return userRole, ok
}
