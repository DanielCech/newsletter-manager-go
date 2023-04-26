package session

import (
	"context"
	"errors"
	"fmt"

	domsession "newsletter-manager-go/domain/session"
	domuser "newsletter-manager-go/domain/user"
	"newsletter-manager-go/types"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/types/id"

	"github.com/prometheus/client_golang/prometheus"
)

// UserService represents object which is capable of reading user in several ways.
type UserService interface {
	Read(ctx context.Context, userID id.User) (*domuser.User, error)
	ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domuser.User, error)
}

// Service consists of session factory and repository and user reader.
type Service struct {
	sessionFactory    domsession.Factory
	sessionRepository domsession.Repository
	userService       UserService
}

// NewService returns new instance of a session service.
func NewService(
	sessionFactory domsession.Factory,
	sessionRepository domsession.Repository,
	userService UserService,
) (*Service, error) {
	if sessionRepository == nil {
		return nil, errors.New("invalid session repository")
	}
	if userService == nil {
		return nil, errors.New("invalid user service")
	}
	return &Service{
		sessionFactory:    sessionFactory,
		sessionRepository: sessionRepository,
		userService:       userService,
	}, nil
}

// Create creates a new session and creates refresh token in the repository.
// Returns a newly created session along with user who is the session owner.
func (s *Service) Create(ctx context.Context, email types.Email, password types.Password) (*domsession.Session, *domuser.User, error) {
	user, err := s.userService.ReadByCredentials(ctx, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("reading user by credentials: %w", err)
	}
	session, err := s.create(ctx, user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	return session, user, nil
}

// CreateForUser creates a new session and creates refresh token in the repository.
// Returns a newly created session.
func (s *Service) CreateForUser(ctx context.Context, user *domuser.User) (*domsession.Session, error) {
	return s.create(ctx, user.ID, user.Role)
}

func (s *Service) create(ctx context.Context, userID id.User, userRole domuser.Role) (*domsession.Session, error) {
	claims, err := domsession.NewClaims(userID, userRole)
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

// DestroyForUser destroys all sessions by deleting refresh tokens from the repository by user id.
func (s *Service) DestroyForUser(ctx context.Context, userID id.User) error {
	if err := s.sessionRepository.DeleteRefreshTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("deleting refresh tokens by user id: %w", err)
	}
	return nil
}

// Refresh reads user which is the token owner and creates a new session.
// Token is then refreshed in the repository.
func (s *Service) Refresh(ctx context.Context, refreshTokenID id.RefreshToken) (session *domsession.Session, err error) {
	err = s.sessionRepository.Refresh(ctx, refreshTokenID, func(oldRefreshToken *domsession.RefreshToken) (*domsession.RefreshToken, error) {
		if oldRefreshToken.IsExpired() {
			return nil, domsession.ErrRefreshTokenExpired
		}
		user, err := s.userService.Read(ctx, oldRefreshToken.UserID)
		if err != nil {
			return nil, fmt.Errorf("reading user: %w", err)
		}
		claims, err := domsession.NewClaims(user.ID, user.Role)
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
