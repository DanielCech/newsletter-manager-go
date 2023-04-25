package context

import (
	"context"

	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"
)

var (
	contextKey = struct {
		userID   ctxKeyUserID
		userRole ctxKeyUserRole
	}{}
)

type (
	ctxKeyUserID   struct{}
	ctxKeyUserRole struct{}
)

// WithUserID passes user ID to the context.
func WithUserID(ctx context.Context, userID id.User) context.Context {
	return context.WithValue(ctx, contextKey.userID, userID)
}

// UserIDFromCtx gets user ID from the context.
func UserIDFromCtx(ctx context.Context) (id.User, bool) {
	userID, ok := ctx.Value(contextKey.userID).(id.User)
	return userID, ok
}

// WithUserRole passes user role to the context.
func WithUserRole(ctx context.Context, role domuser.Role) context.Context {
	return context.WithValue(ctx, contextKey.userRole, role)
}

// UserRoleFromCtx gets user role from the context.
func UserRoleFromCtx(ctx context.Context) (domuser.Role, bool) {
	userRole, ok := ctx.Value(contextKey.userRole).(domuser.Role)
	return userRole, ok
}
