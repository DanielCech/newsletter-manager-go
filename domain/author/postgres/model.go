package author

import (
	"time"

	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/types/id"
)

// author represents table author.
type author struct {
	ID           id.Author `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToAuthor converts author to domain model.
func (u author) ToAuthor(factory domauthor.Factory) *domauthor.Author {
	return factory.NewAuthorFromFields(
		u.ID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
