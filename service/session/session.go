package session

import (
	"context"
	"errors"
	"fmt"

	domauthor "newsletter-manager-go/domain/author"
	domsession "newsletter-manager-go/domain/session"
	"newsletter-manager-go/types"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/types/id"

	"github.com/prometheus/client_golang/prometheus"
)

// authorService represents object which is capable of reading author in several ways.
type authorService interface {
	Read(ctx context.Context, AuthorID id.Author) (*domauthor.Author, error)
	ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domauthor.Author, error)
}

// Service consists of session factory and repository and author reader.
type Service struct {
	sessionFactory    domsession.Factory
	sessionRepository domsession.Repository
	authorService     authorService
}

// NewService returns new instance of a session service.
func NewService(
	sessionFactory domsession.Factory,
	sessionRepository domsession.Repository,
	authorService authorService,
) (*Service, error) {
	if sessionRepository == nil {
		return nil, errors.New("invalid session repository")
	}
	if authorService == nil {
		return nil, errors.New("invalid author service")
	}
	return &Service{
		sessionFactory:    sessionFactory,
		sessionRepository: sessionRepository,
		authorService:     authorService,
	}, nil
}

// Create creates a new session and creates refresh token in the repository.
// Returns a newly created session along with author who is the session owner.
func (s *Service) Create(ctx context.Context, email types.Email, password types.Password) (*domsession.Session, *domauthor.Author, error) {
	author, err := s.authorService.ReadByCredentials(ctx, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("reading author by credentials: %w", err)
	}
	session, err := s.create(ctx, author.ID)
	if err != nil {
		return nil, nil, err
	}
	return session, author, nil
}

// CreateForAuthor creates a new session and creates refresh token in the repository.
// Returns a newly created session.
func (s *Service) CreateForAuthor(ctx context.Context, author *domauthor.Author) (*domsession.Session, error) {
	return s.create(ctx, author.ID)
}

func (s *Service) create(ctx context.Context, authorID id.Author) (*domsession.Session, error) {
	claims, err := domsession.NewClaims(authorID)
	if err != nil {
		return nil, fmt.Errorf("new claims: %w", err)
	}
	session, err := s.sessionFactory.NewSession(claims)
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	if err = s.sessionRepository.CreateRefreshToken(ctx, &session.RefreshToken); err != nil {
		return nil, fmt.Errorf("creating refresh token: %w", err)
	}
	return session, nil
}

// Destroy destroys current session by deleting refresh token from the repository.
func (s *Service) Destroy(ctx context.Context, refreshTokenID id.RefreshToken) error {
	if err := s.sessionRepository.DeleteRefreshToken(ctx, refreshTokenID); err != nil {
		if errors.Is(err, domsession.ErrRefreshTokenNotFound) {
			return apierrors.NewNotFoundError(err, "deleting refresh token").WithPublicMessage(err.Error())
		}
		return fmt.Errorf("deleting refresh token: %w", err)
	}
	return nil
}

// DestroyForAuthor destroys all sessions by deleting refresh tokens from the repository by author id.
func (s *Service) DestroyForAuthor(ctx context.Context, authorID id.Author) error {
	// if err := s.sessionRepository.DeleteRefreshTokensByAuthorID(ctx, AuthorID); err != nil {
	//	return fmt.Errorf("deleting refresh tokens by author id: %w", err)
	// }
	return nil
}

// Refresh reads author which is the token owner and creates a new session.
// Token is then refreshed in the repository.
func (s *Service) Refresh(ctx context.Context, refreshTokenID id.RefreshToken) (session *domsession.Session, err error) {
	err = s.sessionRepository.Refresh(ctx, refreshTokenID, func(oldRefreshToken *domsession.RefreshToken) (*domsession.RefreshToken, error) {
		if oldRefreshToken.IsExpired() {
			return nil, domsession.ErrRefreshTokenExpired
		}
		author, err := s.authorService.Read(ctx, oldRefreshToken.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("reading author: %w", err)
		}
		claims, err := domsession.NewClaims(author.ID)
		if err != nil {
			return nil, fmt.Errorf("new custom claims: %w", err)
		}
		session, err = s.sessionFactory.NewSession(claims)
		if err != nil {
			return nil, fmt.Errorf("new session: %w", err)
		}
		return &session.RefreshToken, nil
	})
	if err != nil {
		if errors.Is(err, domsession.ErrRefreshTokenNotFound) {
			return nil, apierrors.NewNotFoundError(err, "refreshing session").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domsession.ErrRefreshTokenExpired) {
			return nil, apierrors.NewUnauthorizedError(err, "refreshing session").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("refreshing session: %w", err)
	}
	return session, nil
}

func (s *Service) Collect(chan<- prometheus.Metric) {}

func (s *Service) Describe(chan<- *prometheus.Desc) {}
