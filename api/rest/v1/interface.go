package v1

import (
	"context"
	domauthor "newsletter-manager-go/domain/author"
	domsession "newsletter-manager-go/domain/session"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

type AuthorService interface {
	Create(ctx context.Context, CreateAuthorInput domauthor.CreateAuthorInput) (*domauthor.Author, *domsession.Session, error)
	Read(ctx context.Context, authorID id.Author) (*domauthor.Author, error)
	ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error)
	ChangePassword(ctx context.Context, authorID id.Author, oldPassword, newPassword types.Password) error
	List(ctx context.Context) ([]domauthor.Author, error)
}

type NewsletterService interface {
	// Create(ctx context.Context, createNewsletterInput domnewsletter.Newsletter) (*domnewsletter.Newsletter, error)
	// Read(ctx context.Context, NewsletterID id.Newsletter) (*domnewsletter.Newsletter, error)
	// List(ctx context.Context) ([]domnewsletter.Newsletter, error)
}

// SessionService is an interface for v1 session endpoints.
type SessionService interface {
	Create(ctx context.Context, email types.Email, password types.Password) (*domsession.Session, *domauthor.Author, error)
	Destroy(ctx context.Context, refreshTokenID id.RefreshToken) error
	Refresh(ctx context.Context, refreshTokenID id.RefreshToken) (*domsession.Session, error)
}
