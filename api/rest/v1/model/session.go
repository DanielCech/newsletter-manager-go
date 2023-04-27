package model

import (
	"time"

	domsession "newsletter-manager-go/domain/session"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// CreateSessionInput represents JSON body needed for creating a new session.
type CreateSessionInput struct {
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}

// CreateSessionResp represents JSON response body of creating a new session.
type CreateSessionResp struct {
	Author  Author  `json:"user"`
	Session Session `json:"session"`
}

// RefreshSessionInput represents JSON body needed for refreshing an existing session.
type RefreshSessionInput struct {
	RefreshToken id.RefreshToken `json:"refreshToken" validate:"required"`
}

// RefreshSessionResp represents JSON response body of refreshing an existing session.
type RefreshSessionResp struct {
	Session
}

// DestroySessionInput represents JSON body needed for destroying an existing session.
type DestroySessionInput struct {
	RefreshToken id.RefreshToken `json:"refreshToken" validate:"required"`
}

// Session contains access and refresh tokens along with their expiration times.
type Session struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

// FromSession converts domain object to api object.
func FromSession(session *domsession.Session) Session {
	return Session{
		AccessToken:           session.AccessToken.SignedData,
		AccessTokenExpiresAt:  session.AccessToken.ExpiresAt,
		RefreshToken:          string(session.RefreshToken.ID),
		RefreshTokenExpiresAt: session.RefreshToken.ExpiresAt,
	}
}
