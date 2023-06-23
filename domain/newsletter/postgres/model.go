package newsletter

import (
	"time"

	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types/id"
)

// newsletter represents table newsletter.
type newsletter struct {
	ID          id.Newsletter `db:"id"`
	AuthorID    id.Author     `db:"author_id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

// ToNewsletter converts newsletter to domain model.
func (n newsletter) ToDomainNewsletter(factory domnewsletter.Factory) *domnewsletter.Newsletter {
	return factory.NewNewsletterFromFields(
		n.ID,
		n.AuthorID,
		n.Name,
		n.Description,
		n.CreatedAt,
		n.UpdatedAt,
	)
}
