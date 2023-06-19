package model

import (
	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types/id"
)

// Newsletter consists of fields which describe an newsletter.
type Newsletter struct {
	ID          id.Newsletter `json:"id"`
	AuthorID    id.Author     `json:"authorId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
}

// Newsletter consists of fields which describe an newsletter.
type CreateNewsletterReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetNewsletterInput struct {
	AuthorID   id.Author
	Newsletter id.Newsletter
}

// TODO: delete
type NewsletterIDInput struct {
	NewsletterID id.Newsletter `json:"newsletterId"`
}

// FromNewsletter converts domain object to api object.
func FromDomainNewsletter(newsletter *domnewsletter.Newsletter) Newsletter {
	return Newsletter{
		ID:          newsletter.ID,
		AuthorID:    newsletter.AuthorID,
		Name:        newsletter.Name,
		Description: newsletter.Description,
	}
}
