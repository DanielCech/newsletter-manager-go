package newsletter

import (
	"errors"
	domnewsletter "newsletter-manager-go/domain/newsletter"
)

// Service consists of newsletter factory and repository.
type Service struct {
	newsletterFactory    domnewsletter.Factory
	newsletterRepository domnewsletter.Repository
}

// NewService returns new instance of a newsletter service.
func NewService(newsletterFactory domnewsletter.Factory, newsletterRepository domnewsletter.Repository) (*Service, error) {
	if newsletterRepository == nil {
		return nil, errors.New("invalid newsletter repository")
	}
	return &Service{
		newsletterFactory:    newsletterFactory,
		newsletterRepository: newsletterRepository,
	}, nil
}
