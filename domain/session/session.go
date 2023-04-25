package session

import (
	"errors"
	"time"

	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"
	"strv-template-backend-go-api/util/timesource"
)

var (
	ErrInvalidSecret     = errors.New("invalid secret")
	ErrInvalidTimeSource = errors.New("invalid time source")
)

// Claims object contains fields used for access/refresh tokens.
type Claims struct {
	// UserID represents subject claim.
	// See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	UserID id.User
	Custom CustomClaims
}

// NewClaims returns new instance of Claims.
func NewClaims(userID id.User, userRole domuser.Role) (Claims, error) {
	customClaims := Claims{
		UserID: userID,
		Custom: CustomClaims{
			UserRole: userRole,
		},
	}
	if err := customClaims.Valid(); err != nil {
		return Claims{}, err
	}
	return customClaims, nil
}

// Valid checks whether custom jwtClaims are valid.
func (c Claims) Valid() error {
	if c.UserID.Empty() {
		return domuser.ErrInvalidUserID
	}
	return c.Custom.UserRole.Valid()
}

// AccessToken is used for stateless session according to RFC 7519.
type AccessToken struct {
	timeSource timesource.TimeSource

	// Claims object contains session fields.
	Claims Claims
	// SignedData represents signed access token.
	SignedData string
	// ExpiresAt contains access token expiration time.
	ExpiresAt time.Time
}

// IsExpired returns true if access token is already expired.
func (t AccessToken) IsExpired() bool {
	return t.timeSource.Now().After(t.ExpiresAt)
}

// RefreshToken is used for renewal of stateless session.
type RefreshToken struct {
	timeSource timesource.TimeSource

	// ID is safe to return as a refresh token payload.
	ID        id.RefreshToken
	UserID    id.User
	ExpiresAt time.Time
	CreatedAt time.Time
}

// IsExpired returns true if refresh token is already expired.
func (t RefreshToken) IsExpired() bool {
	return t.timeSource.Now().After(t.ExpiresAt)
}

// Session contains access and refresh token.
type Session struct {
	// AccessToken is used for accessing the endpoints.
	AccessToken AccessToken
	// RefreshToken is used for refreshing the session.
	RefreshToken RefreshToken
}

// CustomClaims object contains non-standard fields used for access/refresh tokens.
type CustomClaims struct {
	UserRole domuser.Role `json:"user_role,omitempty"`
}
