package newsletter

import (
	"errors"
	"time"

	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
)

var (
	ErrInvalidNewsletterID       = errors.New("invalid newsletter id")
	ErrInvalidNewsletterName     = errors.New("invalid newsletter name")
	ErrInvalidNewsletterPassword = errors.New("invalid newsletter password")
	ErrInvalidNewsletterRole     = errors.New("invalid newsletter role")
)

// Newsletter consists of fields which describe a newsletter.
type Newsletter struct {
	timeSource timesource.TimeSource

	ID          id.Newsletter
	AuthorID    id.Author
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateAuthorInput consists of fields required for creation of a new author.
type CreateNewsletterInput struct {
	AuthorID    id.Author
	Name        string
	Description string
}

func (u *Newsletter) Valid() error {
	if u.ID.Empty() {
		return ErrInvalidNewsletterID
	}
	if len(u.Name) == 0 {
		return ErrInvalidNewsletterName
	}
	return nil
}
