package model

import (
	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"time"
)

// Author consists of fields which describe an author.
type Author struct {
	ID    id.Author   `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}

type FullAuthor struct {
	ID                     int    `json:"author_id"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	AccessToken            string `json:"access_token,omitempty"`
	AccessTokenExpiration  string `json:"access_token_expiration,omitempty"`
	RefreshToken           string `json:"refresh_token,omitempty"`
	RefreshTokenExpiration string `json:"refresh_token_expiration,omitempty"`
}

// FromAuthor converts domain object to api object.
func FromDomainAuthor(author *domauthor.Author) Author {
	return Author{
		ID:    author.ID,
		Name:  author.Name,
		Email: author.Email,
	}
}

// // FromAuthors converts domain object to api object.
//
//	func FromAuthors(dauthors []domauthor.Author) []Author {
//		authors := make([]Author, 0, len(dauthors))
//		for _, u := range dauthors {
//			authors = append(authors, Author{
//				ID:    u.ID,
//				Name:  u.Name,
//				Email: u.Email,
//			})
//		}
//		return authors
//	}
type CreateAuthorResponse struct {
	Author  *Author `json:"user"`
	Session Session `json:"session"`
}

type AuthorSignUpInput struct {
	Name     string         `json:"name" validate:"required"`
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}

type AuthorSignInInput struct {
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthorIDInput struct {
	AuthorID id.Author `json:"authorId"`
}

type RefreshSessionResponse struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}
