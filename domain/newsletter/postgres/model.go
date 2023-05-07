package newsletter

import (
	"time"

	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types/id"
)

// newsletter represents table newsletter.
type newsletter struct {
	ID           id.Newsletter `db:"id"`
	Name         string        `db:"name"`
	Email        string        `db:"email"`
	PasswordHash []byte        `db:"password_hash"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

// ToNewsletter converts newsletter to domain model.
func (u newsletter) ToDomainNewsletter(factory domnewsletter.Factory) *domnewsletter.Newsletter {
	// TODO
	return nil
}
