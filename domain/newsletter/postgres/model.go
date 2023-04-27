package newsletter

import (
	"time"

	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types/id"
)

// newsletter represents table newsletter.
type newsletter struct {
	ID           id.Newsletter  `db:"id"`
	ReferrerID   *id.Newsletter `db:"referrer_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	PasswordHash []byte         `db:"password_hash"`
	Role         string         `db:"role"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// ToNewsletter converts newsletter to domain model.
func (u newsletter) ToNewsletter(factory domnewsletter.Factory) *domnewsletter.Newsletter {
	return factory.NewNewsletterFromFields(
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
