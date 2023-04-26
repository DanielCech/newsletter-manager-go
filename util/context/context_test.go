package context

import (
	"context"
	"testing"

	domuser "newsletter-manager-go/domain/user"
	"newsletter-manager-go/types/id"

	"github.com/stretchr/testify/assert"
)

func Test_WithUserID(t *testing.T) {
	expected := id.NewUser()
	ctx := WithUserID(context.Background(), expected)
	userID, ok := ctx.Value(contextKey.userID).(id.User)
	assert.True(t, ok)
	assert.Equal(t, expected, userID)
}

func Test_UserIDFromCtx(t *testing.T) {
	expected := id.NewUser()
	ctx := context.WithValue(context.Background(), contextKey.userID, expected)
	userID, ok := UserIDFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, userID)
}

func Test_WithUserRole(t *testing.T) {
	expected := domuser.Role("user")
	ctx := WithUserRole(context.Background(), expected)
	userRole, ok := ctx.Value(contextKey.userRole).(domuser.Role)
	assert.True(t, ok)
	assert.Equal(t, expected, userRole)
}

func Test_UserRoleFromCtx(t *testing.T) {
	expected := domuser.Role("user")
	ctx := context.WithValue(context.Background(), contextKey.userRole, expected)
	userRole, ok := UserRoleFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, userRole)
}
