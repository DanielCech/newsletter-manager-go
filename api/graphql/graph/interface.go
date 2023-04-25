package graph

import (
	"context"

	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"
)

// UserService is an interface for user endpoints.
type UserService interface {
	Create(ctx context.Context, createUserInput domuser.CreateUserInput) (*domuser.User, *domsession.Session, error)
	Read(ctx context.Context, userID id.User) (*domuser.User, error)
	ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domuser.User, error)
	ChangePassword(ctx context.Context, userID id.User, oldPassword, newPassword types.Password) error
	List(ctx context.Context) ([]domuser.User, error)
}

// SessionService is an interface for session endpoints.
type SessionService interface {
	Create(ctx context.Context, email types.Email, password types.Password) (*domsession.Session, *domuser.User, error)
	Destroy(ctx context.Context, refreshTokenID id.RefreshToken) error
	Refresh(ctx context.Context, refreshTokenID id.RefreshToken) (*domsession.Session, error)
}
