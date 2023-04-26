package session

import (
	"context"
	"errors"

	"newsletter-manager-go/types/id"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)

// RefreshFunc is a function for refresh token update.
// It is expected to return a new refresh token based on the existing one.
type RefreshFunc func(*RefreshToken) (*RefreshToken, error)

// Repository consists of functions operating over session objects.
type Repository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) error
	Refresh(ctx context.Context, refreshTokenID id.RefreshToken, fn RefreshFunc) error
	DeleteRefreshToken(ctx context.Context, refreshTokenID id.RefreshToken) error
	DeleteRefreshTokensByUserID(ctx context.Context, userID id.User) error
}
