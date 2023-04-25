package session

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"
	"strv-template-backend-go-api/util/timesource"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

type mockTimeSource struct {
	mock.Mock
}

func (m *mockTimeSource) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateRefreshToken(_ context.Context, refreshToken *domsession.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func refreshTokenMatchFunc(expected *domsession.RefreshToken) func(*domsession.RefreshToken) bool {
	return func(called *domsession.RefreshToken) bool {
		return compareRefreshTokens(expected, called)
	}
}

func (m *mockRepository) Refresh(_ context.Context, id id.RefreshToken, refreshFn domsession.RefreshFunc) error {
	args := m.Called(id, refreshFn)
	return args.Error(0)
}

func refreshMatchFunc(input, expected *domsession.RefreshToken, expectedErr error) func(func(*domsession.RefreshToken) (*domsession.RefreshToken, error)) bool {
	return func(f func(*domsession.RefreshToken) (*domsession.RefreshToken, error)) bool {
		if input == nil {
			return true
		}
		resultSession, err := f(input)
		if !errors.Is(err, expectedErr) {
			return false
		}
		if expectedErr != nil {
			return true
		}
		return compareRefreshTokens(expected, resultSession)
	}
}

func compareRefreshTokens(expected, actual *domsession.RefreshToken) bool {
	return expected.UserID == actual.UserID &&
		expected.ExpiresAt.Equal(actual.ExpiresAt) &&
		expected.CreatedAt.Equal(actual.CreatedAt)
}

func (m *mockRepository) DeleteRefreshToken(_ context.Context, id id.RefreshToken) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockRepository) DeleteRefreshTokensByUserID(_ context.Context, userID id.User) error {
	args := m.Called(userID)
	return args.Error(0)
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Read(_ context.Context, userID id.User) (*domuser.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockUserService) ReadByCredentials(_ context.Context, email types.Email, password types.Password) (*domuser.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*domuser.User), args.Error(1)
}

type mocks struct {
	timeSource  *mockTimeSource
	repository  *mockRepository
	userService *mockUserService
}

func newMocks() mocks {
	return mocks{
		timeSource:  &mockTimeSource{},
		repository:  &mockRepository{},
		userService: &mockUserService{},
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.timeSource.AssertExpectations(t)
	m.repository.AssertExpectations(t)
	m.userService.AssertExpectations(t)
}

func newFactory(t *testing.T, timeSource timesource.TimeSource) domsession.Factory {
	t.Helper()
	sessionFactory, err := domsession.NewFactory(
		[]byte("abc123"),
		timeSource,
		time.Hour,
		time.Hour,
	)
	require.NoError(t, err)
	return sessionFactory
}

func Test_NewService(t *testing.T) {
	mocks := newMocks()
	factory := newFactory(t, mocks.timeSource)
	expected := &Service{
		sessionFactory:    factory,
		sessionRepository: mocks.repository,
		userService:       mocks.userService,
	}
	service, err := NewService(factory, mocks.repository, mocks.userService)
	assert.NoError(t, err)
	assert.Equal(t, expected, service)

	service, err = NewService(factory, mocks.repository, nil)
	assert.EqualError(t, err, "invalid user service")
	assert.Empty(t, service)

	service, err = NewService(factory, nil, nil)
	assert.EqualError(t, err, "invalid session repository")
	assert.Empty(t, service)
}

func Test_Service_Create(t *testing.T) {
	password := types.Password("Topsecret1")
	email := types.Email("jozko.dlouhy@gmail.com")
	user := &domuser.User{
		ID:    id.NewUser(),
		Email: email,
		Role:  domuser.RoleUser,
	}
	customClaims := domsession.Claims{
		UserID: user.ID,
		Custom: domsession.CustomClaims{UserRole: user.Role},
	}
	var session *domsession.Session

	type args struct {
		email    types.Email
		password types.Password
	}
	tests := []struct {
		name            string
		args            args
		mocks           mocks
		expectedSession *domsession.Session
		expectedUser    *domuser.User
		expectedErr     error
	}{
		{
			name: "success",
			args: args{
				email:    email,
				password: password,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ReadByCredentials", email, password).Return(user, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				factory := newFactory(t, mocks.timeSource)
				var err error
				session, err = factory.NewSession(customClaims)
				require.NoError(t, err)
				mocks.repository.On("CreateRefreshToken", mock.MatchedBy(refreshTokenMatchFunc(&session.RefreshToken))).Return(nil)
				return mocks
			}(),
			expectedSession: session,
			expectedUser:    user,
			expectedErr:     nil,
		},
		{
			name: "failure:read-by-credentials",
			args: args{
				email:    email,
				password: password,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ReadByCredentials", email, password).Return((*domuser.User)(nil), errTest)
				return mocks
			}(),
			expectedSession: nil,
			expectedUser:    nil,
			expectedErr:     fmt.Errorf("reading user by credentials: %w", errTest),
		},
		{
			name: "failure:new-custom-claims",
			args: args{
				email:    email,
				password: password,
			},
			mocks: func() mocks {
				mocks := newMocks()
				u := *user
				u.ID = id.User{}
				mocks.userService.On("ReadByCredentials", email, password).Return(&u, nil)
				return mocks
			}(),
			expectedSession: nil,
			expectedUser:    nil,
			expectedErr:     fmt.Errorf("new claims: %w", domuser.ErrInvalidUserID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			session, user, err := service.Create(context.Background(), tt.args.email, tt.args.password)
			if tt.expectedSession != nil {
				assert.Equal(t, tt.expectedSession.AccessToken, session.AccessToken)
				assert.True(t, compareRefreshTokens(&tt.expectedSession.RefreshToken, &session.RefreshToken))
			}
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_CreateForUser(t *testing.T) {
	email := types.Email("jozko.dlouhy@gmail.com")
	user := &domuser.User{
		ID:    id.NewUser(),
		Email: email,
		Role:  domuser.RoleUser,
	}
	customClaims := domsession.Claims{
		UserID: user.ID,
		Custom: domsession.CustomClaims{UserRole: user.Role},
	}
	var session *domsession.Session

	type args struct {
		user *domuser.User
	}
	tests := []struct {
		name            string
		args            args
		mocks           mocks
		expectedSession *domsession.Session
		expectedErr     error
	}{
		{
			name: "success",
			args: args{user: user},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.timeSource.On("Now").Return(time.Now())
				factory := newFactory(t, mocks.timeSource)
				var err error
				session, err = factory.NewSession(customClaims)
				require.NoError(t, err)
				mocks.repository.On("CreateRefreshToken", mock.MatchedBy(refreshTokenMatchFunc(&session.RefreshToken))).Return(nil)
				return mocks
			}(),
			expectedSession: session,
			expectedErr:     nil,
		},
		{
			name:            "failure:new-custom-claims",
			args:            args{user: &domuser.User{}},
			mocks:           newMocks(),
			expectedSession: nil,
			expectedErr:     fmt.Errorf("new claims: %w", domuser.ErrInvalidUserID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			session, err := service.CreateForUser(context.Background(), tt.args.user)
			if tt.expectedSession != nil {
				assert.Equal(t, tt.expectedSession.AccessToken, session.AccessToken)
				assert.True(t, compareRefreshTokens(&tt.expectedSession.RefreshToken, &session.RefreshToken))
			}
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_create(t *testing.T) {
	userID := id.NewUser()
	userRole := domuser.RoleUser
	customClaims := domsession.Claims{
		UserID: userID,
		Custom: domsession.CustomClaims{UserRole: userRole},
	}
	var session *domsession.Session

	type args struct {
		userID   id.User
		userRole domuser.Role
	}
	tests := []struct {
		name            string
		args            args
		mocks           mocks
		expectedSession *domsession.Session
		expectedErr     error
	}{
		{
			name: "success",
			args: args{
				userID:   userID,
				userRole: userRole,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.timeSource.On("Now").Return(time.Now())
				factory := newFactory(t, mocks.timeSource)
				var err error
				session, err = factory.NewSession(customClaims)
				require.NoError(t, err)
				mocks.repository.On("CreateRefreshToken", mock.MatchedBy(refreshTokenMatchFunc(&session.RefreshToken))).Return(nil)
				return mocks
			}(),
			expectedSession: session,
			expectedErr:     nil,
		},
		{
			name: "failure:create-refresh-token",
			args: args{
				userID:   userID,
				userRole: userRole,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.timeSource.On("Now").Return(time.Now())
				factory := newFactory(t, mocks.timeSource)
				refreshToken, err := factory.NewRefreshToken(customClaims.UserID)
				require.NoError(t, err)
				mocks.repository.On("CreateRefreshToken", mock.MatchedBy(refreshTokenMatchFunc(refreshToken))).Return(errTest)
				return mocks
			}(),
			expectedSession: nil,
			expectedErr:     fmt.Errorf("creating refresh token: %w", errTest),
		},
		{
			name:            "failure:new-custom-claims",
			args:            args{},
			mocks:           newMocks(),
			expectedSession: nil,
			expectedErr:     fmt.Errorf("new claims: %w", domuser.ErrInvalidUserID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			session, err := service.create(context.Background(), tt.args.userID, tt.args.userRole)
			if tt.expectedSession != nil {
				assert.Equal(t, tt.expectedSession.AccessToken, session.AccessToken)
				assert.True(t, compareRefreshTokens(&tt.expectedSession.RefreshToken, &session.RefreshToken))
			}
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_Destroy(t *testing.T) {
	refreshTokenID := id.RefreshToken("5asd4a6d4a36d45as36da")

	type args struct {
		refreshTokenID id.RefreshToken
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("DeleteRefreshToken", refreshTokenID).Return(nil)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:delete-refresh-token",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("DeleteRefreshToken", refreshTokenID).Return(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("deleting refresh token: %w", errTest),
		},
		{
			name: "failure:delete-refresh-token/not-found",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("DeleteRefreshToken", refreshTokenID).Return(domsession.ErrRefreshTokenNotFound)
				return mocks
			}(),
			expectedErr: apierrors.NewNotFoundError(domsession.ErrRefreshTokenNotFound, "deleting refresh token").WithPublicMessage(domsession.ErrRefreshTokenNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			err = service.Destroy(context.Background(), tt.args.refreshTokenID)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_DestroyForUser(t *testing.T) {
	userID := id.NewUser()

	type args struct {
		userID id.User
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{userID: userID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("DeleteRefreshTokensByUserID", userID).Return(nil)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:delete-refresh-token",
			args: args{userID: userID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("DeleteRefreshTokensByUserID", userID).Return(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("deleting refresh tokens by user id: %w", errTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			err = service.DestroyForUser(context.Background(), tt.args.userID)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_Refresh(t *testing.T) {
	oldRefreshTokenID := id.RefreshToken("5asd4a6d4a36d45as36da")
	user := &domuser.User{
		ID:   id.User(uuid.MustParse("962b4575-187a-4095-b632-0125527fd4c4")),
		Role: domuser.RoleUser,
	}
	customClaims := domsession.Claims{
		UserID: user.ID,
		Custom: domsession.CustomClaims{UserRole: user.Role},
	}
	var session *domsession.Session

	type args struct {
		refreshTokenID id.RefreshToken
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domsession.Session
		expectedErr error
	}{
		{
			name: "success",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Read", user.ID).Return(user, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				factory := newFactory(t, mocks.timeSource)
				var err error
				session, err = factory.NewSession(customClaims)
				require.NoError(t, err)
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(&session.RefreshToken, &session.RefreshToken, nil))).Return(nil)
				return mocks
			}(),
			expected:    session,
			expectedErr: nil,
		},
		{
			name: "failure:refresh",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(nil, nil, nil))).Return(errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("refreshing session: %w", errTest),
		},
		{
			name: "failure:refresh/token-expired",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				now := time.Now()
				mocks.timeSource.On("Now").Return(now.Add(time.Hour * 24))
				factory := newFactory(t, mocks.timeSource)
				refreshToken := factory.NewRefreshTokenFromFields("5asd4a6d4a36d45as36da", user.ID, now, now)
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(refreshToken, nil, domsession.ErrRefreshTokenExpired))).Return(domsession.ErrRefreshTokenExpired)
				return mocks
			}(),
			expected: nil,
			expectedErr: func() error {
				return apierrors.NewUnauthorizedError(domsession.ErrRefreshTokenExpired, "refreshing session").WithPublicMessage(domsession.ErrRefreshTokenExpired.Error())
			}(),
		},
		{
			name: "failure:refresh/token-not-found",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(nil, nil, nil))).Return(domsession.ErrRefreshTokenNotFound)
				return mocks
			}(),
			expected:    nil,
			expectedErr: apierrors.NewNotFoundError(domsession.ErrRefreshTokenNotFound, "refreshing session").WithPublicMessage(domsession.ErrRefreshTokenNotFound.Error()),
		},
		{
			name: "failure:new-claims",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				now := time.Now()
				mocks.timeSource.On("Now").Return(now)
				factory := newFactory(t, mocks.timeSource)
				refreshToken := factory.NewRefreshTokenFromFields("5asd4a6d4a36d45as36da", user.ID, now.Add(time.Hour), now)
				u := *user
				u.Role = "unknown"
				mocks.userService.On("Read", u.ID).Return(&u, nil)
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(refreshToken, nil, domuser.ErrInvalidUserRole))).Return(fmt.Errorf("new custom claims: %w", domuser.ErrInvalidUserRole))
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("refreshing session: %w", fmt.Errorf("new custom claims: %w", domuser.ErrInvalidUserRole)),
		},
		{
			name: "failure:read-user",
			args: args{refreshTokenID: oldRefreshTokenID},
			mocks: func() mocks {
				mocks := newMocks()
				now := time.Now()
				mocks.timeSource.On("Now").Return(now)
				factory := newFactory(t, mocks.timeSource)
				refreshToken := factory.NewRefreshTokenFromFields("5asd4a6d4a36d45as36da", user.ID, now.Add(time.Hour), now)
				mocks.userService.On("Read", user.ID).Return((*domuser.User)(nil), errTest)
				mocks.repository.On("Refresh", oldRefreshTokenID, mock.MatchedBy(refreshMatchFunc(refreshToken, nil, errTest))).Return(fmt.Errorf("reading user: %w", errTest))
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("refreshing session: %w", fmt.Errorf("reading user: %w", errTest)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.userService)
			assert.NoError(t, err)
			session, err := service.Refresh(context.Background(), tt.args.refreshTokenID)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expected != nil {
				assert.Equal(t, tt.expected.AccessToken.Claims, session.AccessToken.Claims)
				assert.True(t, compareRefreshTokens(&tt.expected.RefreshToken, &session.RefreshToken))
			}
			tt.mocks.assertExpectations(t)
		})
	}
}
