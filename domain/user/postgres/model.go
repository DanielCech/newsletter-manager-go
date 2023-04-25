package user

import (
	"time"

	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"
)

// user represents table user.
type user struct {
	ID           id.User   `db:"id"`
	ReferrerID   *id.User  `db:"referrer_id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToUser converts user to domain model.
func (u user) ToUser(factory domuser.Factory) *domuser.User {
	return factory.NewUserFromFields(
		u.ID,
		u.ReferrerID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
