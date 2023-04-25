package graph

import (
	"context"
	"testing"

	"strv-template-backend-go-api/api/graphql/graph/model"
	domuser "strv-template-backend-go-api/domain/user"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"
	utilctx "strv-template-backend-go-api/util/context"

	"github.com/stretchr/testify/assert"
)

func Test_Directive_Auth(t *testing.T) {
	directive := NewDirectiveHandler(true)
	next := func(ctx context.Context) (any, error) {
		return nil, nil
	}

	ctx := utilctx.WithUserID(context.Background(), id.NewUser())
	ctx = utilctx.WithUserRole(ctx, domuser.RoleAdmin)
	result, err := directive.Auth(ctx, false, next, []model.Role{model.RoleAdmin})
	assert.NoError(t, err)
	assert.Nil(t, result)

	ctx = utilctx.WithUserID(context.Background(), id.NewUser())
	ctx = utilctx.WithUserRole(ctx, domuser.RoleUser)
	result, err = directive.Auth(ctx, false, next, []model.Role{model.RoleAdmin})
	var e *apierrors.Error
	assert.ErrorAs(t, err, &e)
	assert.Equal(t, apierrors.CodeForbidden, e.Code)
	assert.Equal(t, "insufficient user role", e.PublicMessage)
	assert.Nil(t, result)

	ctx = utilctx.WithUserID(context.Background(), id.NewUser())
	result, err = directive.Auth(ctx, false, next, nil)
	assert.EqualError(t, err, "missing user role in context")
	assert.Nil(t, result)

	result, err = directive.Auth(context.Background(), false, next, nil)
	assert.EqualError(t, err, "missing user ID in context")
	assert.Nil(t, result)

	directive = NewDirectiveHandler(false)
	next = func(ctx context.Context) (any, error) {
		return nil, nil
	}
	result, err = directive.Auth(context.Background(), false, next, nil)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func Test_isSufficientRole(t *testing.T) {
	assert.True(t, isSufficientRole(domuser.RoleAdmin, []model.Role{model.RoleAdmin}))
	assert.False(t, isSufficientRole(domuser.RoleUser, []model.Role{model.RoleAdmin}))
	assert.False(t, isSufficientRole(domuser.RoleAdmin, []model.Role{"invalid"}))
}
