package author

import (
	"context"
	"errors"
	"fmt"
	domsession "newsletter-manager-go/domain/session"

	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/types"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/types/id"

	"github.com/prometheus/client_golang/prometheus"
)

// SessionService represents object which is capable of:
//   - Creating new session
//   - Destroying all author's sessions
type SessionService interface {
	CreateForAuthor(ctx context.Context, author *domauthor.Author) (*domsession.Session, error)
	DestroyForAuthor(ctx context.Context, authorID id.Author) error
}

// Service consists of author factory and repository.
type Service struct {
	authorFactory    domauthor.Factory
	authorRepository domauthor.Repository
	sessionService   SessionService
}

// NewService returns new instance of a author service.
func NewService(authorFactory domauthor.Factory, authorRepository domauthor.Repository, sessionCreator SessionService) (*Service, error) {
	if authorRepository == nil {
		return nil, errors.New("invalid author repository")
	}
	if sessionCreator == nil {
		return nil, errors.New("invalid session service")
	}
	return &Service{
		authorFactory:    authorFactory,
		authorRepository: authorRepository,
		sessionService:   sessionCreator,
	}, nil
}

// Create creates a new author and creates him in the repository.
func (s *Service) Create(ctx context.Context, createAuthorInput domauthor.CreateAuthorInput) (*domauthor.Author, *domsession.Session, error) {
	author, err := s.authorFactory.NewAuthor(createAuthorInput, domauthor.RoleAuthor)
	if err != nil {
		return nil, nil, fmt.Errorf("new author: %w", err)
	}
	if err = s.authorRepository.Create(ctx, author); err != nil {
		if errors.Is(err, domauthor.ErrAuthorEmailAlreadyExists) {
			return nil, nil, apierrors.NewAlreadyExistsError(err, "creating author").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domauthor.ErrReferrerNotFound) {
			return nil, nil, apierrors.NewBadRequestError(err, "creating author").WithPublicMessage(err.Error())
		}
		return nil, nil, fmt.Errorf("creating author: %w", err)
	}
	session, err := s.sessionService.CreateForAuthor(ctx, author)
	if err != nil {
		return nil, nil, fmt.Errorf("creating session for author: %w", err)
	}
	return author, session, nil
}

// Read reads an existing author from the repository.
func (s *Service) Read(ctx context.Context, authorID id.Author) (*domauthor.Author, error) {
	author, err := s.authorRepository.Read(ctx, authorID)
	if err != nil {
		if errors.Is(err, domauthor.ErrAuthorNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading author").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading author: %w", err)
	}
	return author, nil
}

// ReadByEmail reads an existing author from the repository by email.
func (s *Service) ReadByEmail(ctx context.Context, email types.Email) (*domauthor.Author, error) {
	author, err := s.authorRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domauthor.ErrAuthorNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading author by email").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading author by email: %w", err)
	}
	return author, nil
}

// ReadByCredentials reads an existing author from the repository by credentials.
func (s *Service) ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error) {
	const publicErrMsg = "email or password is incorrect"
	author, err := s.authorRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domauthor.ErrAuthorNotFound) {
			return nil, apierrors.NewUnauthorizedError(err, "reading author by credentials").WithPublicMessage(publicErrMsg)
		}
		return nil, fmt.Errorf("reading author by credentials: %w", err)
	}

	if !author.MatchPassword(password) {
		err = errors.New("invalid password")
		return nil, apierrors.NewUnauthorizedError(err, "").WithPublicMessage(publicErrMsg)
	}

	return author, nil
}

// ChangePassword changes author password and updates author in the repository.
func (s *Service) ChangePassword(ctx context.Context, authorID id.Author, oldPassword, newPassword types.Password) error {
	err := s.authorRepository.Update(ctx, authorID, func(u *domauthor.Author) (*domauthor.Author, error) {
		if err := u.ChangePassword(oldPassword, newPassword); err != nil {
			return nil, err
		}
		return u, nil
	})
	if err != nil {
		if errors.Is(err, domauthor.ErrAuthorNotFound) {
			return apierrors.NewUnauthorizedError(err, "changing password").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domauthor.ErrInvalidAuthorPassword) {
			return apierrors.NewBadRequestError(err, "changing password").WithPublicMessage(err.Error())
		}
		return fmt.Errorf("changing password: %w", err)
	}
	if err = s.sessionService.DestroyForAuthor(ctx, authorID); err != nil {
		return fmt.Errorf("destroying sessions for author: %w", err)
	}
	return nil
}

// List lists authors from repository.
func (s *Service) List(ctx context.Context) ([]domauthor.Author, error) {
	authors, err := s.authorRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing authors: %w", err)
	}
	return authors, nil
}

func (s *Service) Collect(chan<- prometheus.Metric) {}

func (s *Service) Describe(chan<- *prometheus.Desc) {}
