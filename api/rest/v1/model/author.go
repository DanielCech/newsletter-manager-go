package model

import (
	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Author consists of fields which describe an author.
type Author struct {
	ID    id.Author   `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}

// FromAuthor converts domain object to api object.
func FromAuthor(author *domauthor.Author) Author {
	return Author{
		ID:    author.ID,
		Name:  author.Name,
		Email: author.Email,
	}
}

// FromAuthors converts domain object to api object.
func FromAuthors(dauthors []domauthor.Author) []Author {
	authors := make([]Author, 0, len(dauthors))
	for _, u := range dauthors {
		authors = append(authors, Author{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return authors
}

// CreateAuthorInput represents JSON body needed for creating a new author.
type CreateAuthorInput struct {
	Name       string         `json:"name" validate:"required"`
	Email      types.Email    `json:"email"`
	Password   types.Password `json:"password"`
	ReferrerID *id.Author     `json:"referrerId"`
}

// CreateAuthorResp represents JSON response body of creating a new author.
type CreateAuthorResp struct {
	Author  Author  `json:"author"`
	Session Session `json:"session"`
}

// ChangeAuthorPasswordInput represents JSON body needed for changing the author password.
type ChangeAuthorPasswordInput struct {
	OldPassword types.Password `json:"oldPassword"`
	NewPassword types.Password `json:"newPassword"`
}
