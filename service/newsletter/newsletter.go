package newsletter

import (
	"context"
	"errors"
	"fmt"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/types/id"
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

func (s *Service) Create(ctx context.Context, createNewsletterInput domnewsletter.CreateNewsletterInput) (*domnewsletter.Newsletter, error) {

	newsletter, err := s.newsletterFactory.NewNewsletter(createNewsletterInput)
	if err != nil {
		return nil, fmt.Errorf("creating newsletter: %w", err)
	}

	err = s.newsletterRepository.Create(ctx, newsletter)

	if err != nil {
		return nil, fmt.Errorf("creating newsletter: %w", err)
	}

	return newsletter, nil
}

func (s *Service) Read(ctx context.Context, newsletterID id.Newsletter) (*domnewsletter.Newsletter, error) {
	newsletter, err := s.newsletterRepository.Read(ctx, newsletterID)
	if err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading newsletter")
		}
		return nil, fmt.Errorf("reading event: %w", err)
	}
	return newsletter, nil
}

// // List lists authors from repository.
func (s *Service) ListCurrentAuthorNewsletters(ctx context.Context, authorID id.Author) ([]domnewsletter.Newsletter, error) {
	return nil, nil
	// authors, err := s.newsletterRepository.ListByAuthor(ctx, authorID)
	// if err != nil {
	// 	return nil, fmt.Errorf("listing author's newsletters: %w", err)
	// }
	// return authors, nil
}
