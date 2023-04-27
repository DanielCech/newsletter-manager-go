package v1

import (
	"context"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	domsession "newsletter-manager-go/domain/session"

	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// AuthorService is an interface for v1 user endpoints.
type AuthorService interface {
	Create(ctx context.Context, createUserInput domauthor.CreateAuthorInput) (*domauthor.Author, *domsession.Session, error)
	Read(ctx context.Context, AuthorID id.Author) (*domauthor.Author, error)
	ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error)
	ChangePassword(ctx context.Context, AuthorID id.Author, oldPassword, newPassword types.Password) error
	List(ctx context.Context) ([]domauthor.Author, error)
}

// NewsletterService is an interface for v1 user endpoints.
type NewsletterService interface {
	Create(ctx context.Context, createNewsletterInput domnewsletter.CreateNewsletterInput) (*domnewsletter.Newsletter, *domsession.Session, error)
	Read(ctx context.Context, NewsletterID id.Newsletter) (*domnewsletter.Newsletter, error)
	List(ctx context.Context) ([]domnewsletter.Newsletter, error)
}

// SessionService is an interface for v1 session endpoints.
type SessionService interface {
	Create(ctx context.Context, email types.Email, password types.Password) (*domsession.Session, *domauthor.Author, error)
	Destroy(ctx context.Context, refreshTokenID id.RefreshToken) error
	Refresh(ctx context.Context, refreshTokenID id.RefreshToken) (*domsession.Session, error)
}
