package newsletter

import (
	"context"
	"errors"
	"fmt"
	domsession "newsletter-manager-go/domain/session"

	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/types/id"

	"github.com/prometheus/client_golang/prometheus"
)

// SessionService represents object which is capable of:
//   - Creating new session
//   - Destroying all newsletter's sessions
type SessionService interface {
	CreateForNewsletter(ctx context.Context, newsletter *domnewsletter.Newsletter) (*domsession.Session, error)
	DestroyForNewsletter(ctx context.Context, newsletterID id.Newsletter) error
}

// Service consists of newsletter factory and repository.
type Service struct {
	newsletterFactory    domnewsletter.Factory
	newsletterRepository domnewsletter.Repository
	sessionService       SessionService
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

// Create creates a new newsletter and creates him in the repository.
func (s *Service) Create(ctx context.Context, createNewsletterInput domnewsletter.CreateNewsletterInput) (*domnewsletter.Newsletter, *domsession.Session, error) {
	newsletter, err := s.newsletterFactory.NewNewsletter(createNewsletterInput, domnewsletter.RoleNewsletter)
	if err != nil {
		return nil, nil, fmt.Errorf("new newsletter: %w", err)
	}
	if err = s.newsletterRepository.Create(ctx, newsletter); err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterEmailAlreadyExists) {
			return nil, nil, apierrors.NewAlreadyExistsError(err, "creating newsletter").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domnewsletter.ErrReferrerNotFound) {
			return nil, nil, apierrors.NewBadRequestError(err, "creating newsletter").WithPublicMessage(err.Error())
		}
		return nil, nil, fmt.Errorf("creating newsletter: %w", err)
	}
	session, err := s.sessionService.CreateForNewsletter(ctx, newsletter)
	if err != nil {
		return nil, nil, fmt.Errorf("creating session for newsletter: %w", err)
	}
	return newsletter, session, nil
}

// Read reads an existing newsletter from the repository.
func (s *Service) Read(ctx context.Context, newsletterID id.Newsletter) (*domnewsletter.Newsletter, error) {
	newsletter, err := s.newsletterRepository.Read(ctx, newsletterID)
	if err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading newsletter").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading newsletter: %w", err)
	}
	return newsletter, nil
}

// ReadByEmail reads an existing newsletter from the repository by email.
func (s *Service) ReadByEmail(ctx context.Context, email types.Email) (*domnewsletter.Newsletter, error) {
	newsletter, err := s.newsletterRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading newsletter by email").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading newsletter by email: %w", err)
	}
	return newsletter, nil
}

// ReadByCredentials reads an existing newsletter from the repository by credentials.
func (s *Service) ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domnewsletter.Newsletter, error) {
	const publicErrMsg = "email or password is incorrect"
	newsletter, err := s.newsletterRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterNotFound) {
			return nil, apierrors.NewUnauthorizedError(err, "reading newsletter by credentials").WithPublicMessage(publicErrMsg)
		}
		return nil, fmt.Errorf("reading newsletter by credentials: %w", err)
	}

	if !newsletter.MatchPassword(password) {
		err = errors.New("invalid password")
		return nil, apierrors.NewUnauthorizedError(err, "").WithPublicMessage(publicErrMsg)
	}

	return newsletter, nil
}

// ChangePassword changes newsletter password and updates newsletter in the repository.
func (s *Service) ChangePassword(ctx context.Context, newsletterID id.Newsletter, oldPassword, newPassword types.Password) error {
	err := s.newsletterRepository.Update(ctx, newsletterID, func(u *domnewsletter.Newsletter) (*domnewsletter.Newsletter, error) {
		if err := u.ChangePassword(oldPassword, newPassword); err != nil {
			return nil, err
		}
		return u, nil
	})
	if err != nil {
		if errors.Is(err, domnewsletter.ErrNewsletterNotFound) {
			return apierrors.NewUnauthorizedError(err, "changing password").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domnewsletter.ErrInvalidNewsletterPassword) {
			return apierrors.NewBadRequestError(err, "changing password").WithPublicMessage(err.Error())
		}
		return fmt.Errorf("changing password: %w", err)
	}
	if err = s.sessionService.DestroyForNewsletter(ctx, newsletterID); err != nil {
		return fmt.Errorf("destroying sessions for newsletter: %w", err)
	}
	return nil
}

// List lists newsletters from repository.
func (s *Service) List(ctx context.Context) ([]domnewsletter.Newsletter, error) {
	newsletters, err := s.newsletterRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing newsletters: %w", err)
	}
	return newsletters, nil
}

func (s *Service) Collect(chan<- prometheus.Metric) {}

func (s *Service) Describe(chan<- *prometheus.Desc) {}
