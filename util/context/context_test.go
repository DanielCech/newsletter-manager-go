package context

import (
	"context"
	"testing"

	domauthor "newsletter-manager-go/domain/user"
	"newsletter-manager-go/types/id"

	"github.com/stretchr/testify/assert"
)

func Test_WithAuthorID(t *testing.T) {
	expected := id.NewUser()
	ctx := WithAuthorID(context.Background(), expected)
	AuthorID, ok := ctx.Value(contextKey.AuthorID).(id.Author)
	assert.True(t, ok)
	assert.Equal(t, expected, AuthorID)
}

func Test_AuthorIDFromCtx(t *testing.T) {
	expected := id.NewUser()
	ctx := context.WithValue(context.Background(), contextKey.AuthorID, expected)
	AuthorID, ok := AuthorIDFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, AuthorID)
}

func Test_WithUserRole(t *testing.T) {
	expected := domauthor.Role("user")
	ctx := WithUserRole(context.Background(), expected)
	userRole, ok := ctx.Value(contextKey.userRole).(domauthor.Role)
	assert.True(t, ok)
	assert.Equal(t, expected, userRole)
}

func Test_UserRoleFromCtx(t *testing.T) {
	expected := domauthor.Role("user")
	ctx := context.WithValue(context.Background(), contextKey.userRole, expected)
	userRole, ok := UserRoleFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, userRole)
}
