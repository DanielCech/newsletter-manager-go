package graph

import (
	"context"
	"errors"
	"strings"

	"strv-template-backend-go-api/api/graphql/graph/model"
	domuser "strv-template-backend-go-api/domain/user"
	apierrors "strv-template-backend-go-api/types/errors"
	utilctx "strv-template-backend-go-api/util/context"

	"github.com/99designs/gqlgen/graphql"
)

// NewDirectiveHandler creates a new DirectiveRoot.
func NewDirectiveHandler(authRequired bool) DirectiveRoot {
	d := directive{authRequired: authRequired}
	return DirectiveRoot{
		Auth: d.Auth,
	}
}

type directive struct {
	authRequired bool
}

// Auth verifies a user has the correct permissions to access an operation.
func (d *directive) Auth(ctx context.Context, _ any, next graphql.Resolver, sufficientRoles []model.Role) (any, error) {
	if !d.authRequired {
		return next(ctx)
	}

	_, ok := utilctx.UserIDFromCtx(ctx)
	if !ok {
		return nil, errors.New("missing user ID in context")
	}

	userRole, ok := utilctx.UserRoleFromCtx(ctx)
	if !ok {
		return nil, errors.New("missing user role in context")
	}

	if !isSufficientRole(userRole, sufficientRoles) {
		err := errors.New("insufficient user role")
		return nil, apierrors.NewForbiddenError(err, "").WithPublicMessage(err.Error())
	}

	return next(ctx)
}

func isSufficientRole(role domuser.Role, sufficientRoles []model.Role) bool {
	for _, sufficientRole := range sufficientRoles {
		domUserRole, err := domuser.NewRole(strings.ToLower(sufficientRole.String()))
		if err != nil {
			return false
		}
		if role.IsSufficientToRole(domUserRole) {
			return true
		}
	}
	return false
}
