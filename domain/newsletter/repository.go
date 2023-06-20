package newsletter

import (
	"context"
	"errors"

	"newsletter-manager-go/types/id"
)

var (
	ErrNewsletterEmailAlreadyExists = errors.New("newsletter email already exists")
	ErrNewsletterNotFound           = errors.New("newsletter not found")
	ErrReferrerNotFound             = errors.New("referrer not found")
)

// UpdateFunc is a function for newsletter update.
// It is expected to return a new newsletter based on the existing one.
type UpdateFunc func(*Newsletter) (*Newsletter, error)

// Repository consists of functions operating over newsletter objects.
type Repository interface {
	Create(ctx context.Context, newsletter *Newsletter) error
	Read(ctx context.Context, newsletterID id.Newsletter) (*Newsletter, error)
	// ReadByEmail(ctx context.Context, email types.Email) (*Newsletter, error)
	// List(ctx context.Context) ([]Newsletter, error)
	// ListByAuthor(ctx context.Context, authorID id.Author) ([]Newsletter, error)
	Update(ctx context.Context, newsletterID id.Newsletter, fn UpdateFunc) error
}
