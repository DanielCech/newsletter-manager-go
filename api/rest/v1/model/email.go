package model

import (
	domemail "newsletter-manager-go/domain/email"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Email consists of fields which describe an email.
type Email struct {
	ID    id.Email    `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}

// FromEmail converts domain object to api object.
func FromEmail(email *domemail.Email) Email {
	return Email{
		ID:    email.ID,
		Name:  email.Name,
		Email: email.Email,
	}
}

// FromEmails converts domain object to api object.
func FromEmails(demails []domemail.Email) []Email {
	emails := make([]Email, 0, len(demails))
	for _, u := range demails {
		emails = append(emails, Email{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return emails
}

// CreateEmailInput represents JSON body needed for creating a new email.
type CreateEmailInput struct {
	Name     string         `json:"name" validate:"required"`
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}

// CreateEmailResp represents JSON response body of creating a new email.
type CreateEmailResp struct {
	Email   Email   `json:"email"`
	Session Session `json:"session"`
}
