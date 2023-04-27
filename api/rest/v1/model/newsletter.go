package model

import (
	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Newsletter consists of fields which describe an newsletter.
type Newsletter struct {
	ID    id.Newsletter `json:"id"`
	Name  string        `json:"name"`
	Email types.Email   `json:"email"`
}

// FromNewsletter converts domain object to api object.
func FromNewsletter(newsletter *domnewsletter.Newsletter) Newsletter {
	return Newsletter{
		ID:    newsletter.ID,
		Name:  newsletter.Name,
		Email: newsletter.Email,
	}
}

// FromNewsletters converts domain object to api object.
func FromNewsletters(dnewsletters []domnewsletter.Newsletter) []Newsletter {
	newsletters := make([]Newsletter, 0, len(dnewsletters))
	for _, u := range dnewsletters {
		newsletters = append(newsletters, Newsletter{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return newsletters
}

// CreateNewsletterInput represents JSON body needed for creating a new newsletter.
type CreateNewsletterInput struct {
	Name       string         `json:"name" validate:"required"`
	Email      types.Email    `json:"email"`
	Password   types.Password `json:"password"`
	ReferrerID *id.Newsletter `json:"referrerId"`
}

// CreateNewsletterResp represents JSON response body of creating a new newsletter.
type CreateNewsletterResp struct {
	Newsletter Newsletter `json:"newsletter"`
	Session    Session    `json:"session"`
}
