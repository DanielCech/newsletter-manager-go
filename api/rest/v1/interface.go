package v1

import (
	"context"
	domauthor "newsletter-manager-go/domain/author"
)

type AuthorService interface {
	Create(ctx context.Context, createUserInput domauthor.CreateAuthorInput) (*domauthor.Author, error)
	//Read(ctx context.Context, AuthorID id.Author) (*domauthor.Author, error)
	//ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error)
	//ChangePassword(ctx context.Context, AuthorID id.Author, oldPassword, newPassword types.Password) error
	//List(ctx context.Context) ([]domauthor.Author, error)
}

type NewsletterService interface {
	//Create(ctx context.Context, createNewsletterInput domnewsletter.Newsletter) (*domnewsletter.Newsletter, error)
	//Read(ctx context.Context, NewsletterID id.Newsletter) (*domnewsletter.Newsletter, error)
	//List(ctx context.Context) ([]domnewsletter.Newsletter, error)
}

//type SessionService interface {
//	Create(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error)
//	Destroy(ctx context.Context, refreshTokenID id.RefreshToken) error
//	Refresh(ctx context.Context, refreshTokenID id.RefreshToken) error
//}
