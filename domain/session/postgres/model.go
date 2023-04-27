package session

import (
	"time"

	domsession "strv-template-backend-go-api/domain/session"
	"strv-template-backend-go-api/types/id"
)

// refreshToken represents table refresh_token.
type refreshToken struct {
	ID        id.RefreshToken `db:"id"`
	AuthorID  id.Author       `db:"user_id"`
	UserRole  string          `db:"user_role"`
	ExpiresAt time.Time       `db:"expires_at"`
	CreatedAt time.Time       `db:"created_at"`
}

// ToRefreshToken converts refreshToken to domain model.
func (r refreshToken) ToRefreshToken(factory domsession.Factory) *domsession.RefreshToken {
	return factory.NewRefreshTokenFromFields(
		r.ID,
		r.AuthorID,
		r.ExpiresAt,
		r.CreatedAt,
	)
}
