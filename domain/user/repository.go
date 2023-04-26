package user

import (
	"context"
	"errors"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

var (
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrReferrerNotFound       = errors.New("referrer not found")
)

// UpdateFunc is a function for user update.
// It is expected to return a new user based on the existing one.
type UpdateFunc func(*User) (*User, error)

// Repository consists of functions operating over user objects.
type Repository interface {
	Create(ctx context.Context, user *User) error
	Read(ctx context.Context, userID id.User) (*User, error)
	ReadByEmail(ctx context.Context, email types.Email) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, userID id.User, fn UpdateFunc) error
}
