package author

import (
	"context"
	"errors"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

var (
	ErrAuthorEmailAlreadyExists = errors.New("author email already exists")
	ErrAuthorNotFound           = errors.New("author not found")
	ErrReferrerNotFound         = errors.New("referrer not found")
)

// UpdateFunc is a function for author update.
// It is expected to return a new author based on the existing one.
type UpdateFunc func(*Author) (*Author, error)

// Repository consists of functions operating over author objects.
type Repository interface {
	Create(ctx context.Context, author *Author) error
	Read(ctx context.Context, authorID id.Author) (*Author, error)
	ReadByEmail(ctx context.Context, email types.Email) (*Author, error)
	List(ctx context.Context) ([]Author, error)
	Update(ctx context.Context, authorID id.Author, fn UpdateFunc) error
	Delete(ctx context.Context, authorID id.Author) error
}
