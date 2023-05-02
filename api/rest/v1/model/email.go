package model

import (
	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"time"
)

// Email consists of fields which describe an email.
type Email struct {
	ID    id.Email    `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}

type CreateEmailInput struct {
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}

type FullEmail struct {
	ID           id.Email      `json:"id"`
	NewsletterID id.Newsletter `json:"newsletterIdd"`
	Subject      string        `json:"subject"`
	HtmlContent  string        `json:"htmlContent"`
	Date         time.Time     `json:"date"`
}

// FromEmail converts domain object to api object.
func FromDomainEmail(email *domnewsletter.Email) Email {
	return Email{
		ID:    email.ID,
		Name:  email.Name,
		Email: email.Email,
	}
}
