package user

import (
	"context"
	"errors"
	"fmt"

	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"

	"github.com/prometheus/client_golang/prometheus"
)

// SessionService represents object which is capable of:
//   - Creating new session
//   - Destroying all user's sessions
type SessionService interface {
	CreateForUser(ctx context.Context, user *domuser.User) (*domsession.Session, error)
	DestroyForUser(ctx context.Context, userID id.User) error
}

// Service consists of user factory and repository.
type Service struct {
	userFactory    domuser.Factory
	userRepository domuser.Repository
	sessionService SessionService
}

// NewService returns new instance of a user service.
func NewService(userFactory domuser.Factory, userRepository domuser.Repository, sessionCreator SessionService) (*Service, error) {
	if userRepository == nil {
		return nil, errors.New("invalid user repository")
	}
	if sessionCreator == nil {
		return nil, errors.New("invalid session service")
	}
	return &Service{
		userFactory:    userFactory,
		userRepository: userRepository,
		sessionService: sessionCreator,
	}, nil
}

// Create creates a new user and creates him in the repository.
func (s *Service) Create(ctx context.Context, createUserInput domuser.CreateUserInput) (*domuser.User, *domsession.Session, error) {
	user, err := s.userFactory.NewUser(createUserInput, domuser.RoleUser)
	if err != nil {
		return nil, nil, fmt.Errorf("new user: %w", err)
	}
	if err = s.userRepository.Create(ctx, user); err != nil {
		if errors.Is(err, domuser.ErrUserEmailAlreadyExists) {
			return nil, nil, apierrors.NewAlreadyExistsError(err, "creating user").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domuser.ErrReferrerNotFound) {
			return nil, nil, apierrors.NewBadRequestError(err, "creating user").WithPublicMessage(err.Error())
		}
		return nil, nil, fmt.Errorf("creating user: %w", err)
	}
	session, err := s.sessionService.CreateForUser(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("creating session for user: %w", err)
	}
	return user, session, nil
}

// Read reads an existing user from the repository.
func (s *Service) Read(ctx context.Context, userID id.User) (*domuser.User, error) {
	user, err := s.userRepository.Read(ctx, userID)
	if err != nil {
		if errors.Is(err, domuser.ErrUserNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading user").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading user: %w", err)
	}
	return user, nil
}

// ReadByEmail reads an existing user from the repository by email.
func (s *Service) ReadByEmail(ctx context.Context, email types.Email) (*domuser.User, error) {
	user, err := s.userRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domuser.ErrUserNotFound) {
			return nil, apierrors.NewNotFoundError(err, "reading user by email").WithPublicMessage(err.Error())
		}
		return nil, fmt.Errorf("reading user by email: %w", err)
	}
	return user, nil
}

// ReadByCredentials reads an existing user from the repository by credentials.
func (s *Service) ReadByCredentials(ctx context.Context, email types.Email, password types.Password) (*domuser.User, error) {
	const publicErrMsg = "email or password is incorrect"
	user, err := s.userRepository.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domuser.ErrUserNotFound) {
			return nil, apierrors.NewUnauthorizedError(err, "reading user by credentials").WithPublicMessage(publicErrMsg)
		}
		return nil, fmt.Errorf("reading user by credentials: %w", err)
	}

	if !user.MatchPassword(password) {
		err = errors.New("invalid password")
		return nil, apierrors.NewUnauthorizedError(err, "").WithPublicMessage(publicErrMsg)
	}

	return user, nil
}

// ChangePassword changes user password and updates user in the repository.
func (s *Service) ChangePassword(ctx context.Context, userID id.User, oldPassword, newPassword types.Password) error {
	err := s.userRepository.Update(ctx, userID, func(u *domuser.User) (*domuser.User, error) {
		if err := u.ChangePassword(oldPassword, newPassword); err != nil {
			return nil, err
		}
		return u, nil
	})
	if err != nil {
		if errors.Is(err, domuser.ErrUserNotFound) {
			return apierrors.NewUnauthorizedError(err, "changing password").WithPublicMessage(err.Error())
		}
		if errors.Is(err, domuser.ErrInvalidUserPassword) {
			return apierrors.NewBadRequestError(err, "changing password").WithPublicMessage(err.Error())
		}
		return fmt.Errorf("changing password: %w", err)
	}
	if err = s.sessionService.DestroyForUser(ctx, userID); err != nil {
		return fmt.Errorf("destroying sessions for user: %w", err)
	}
	return nil
}

// List lists users from repository.
func (s *Service) List(ctx context.Context) ([]domuser.User, error) {
	users, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return users, nil
}

func (s *Service) Collect(chan<- prometheus.Metric) {}

func (s *Service) Describe(chan<- *prometheus.Desc) {}
