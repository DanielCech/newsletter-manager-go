package user

import (
	"bytes"
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testUserName     = "Jozko Dlouhy"
	testUserEmail    = types.Email("jozko.dlouhy@gmail.com")
	testUserPassword = types.Password("Topsecret1")
)

var (
	errTest = errors.New("test error")
)

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) HashPassword(password []byte) ([]byte, error) {
	args := m.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockHasher) CompareHashAndPassword(hash, password []byte) bool {
	args := m.Called(hash, password)
	return args.Bool(0)
}

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

func (m *mockRepository) Create(_ context.Context, user *domuser.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func createUserMatchFunc(expected *domuser.User) func(*domuser.User) bool {
	return func(called *domuser.User) bool {
		return compareUsers(expected, called)
	}
}

func (m *mockRepository) Read(_ context.Context, id id.User) (*domuser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockRepository) ReadByEmail(_ context.Context, email types.Email) (*domuser.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockRepository) List(context.Context) ([]domuser.User, error) {
	args := m.Called()
	return args.Get(0).([]domuser.User), args.Error(1)
}

func (m *mockRepository) Update(_ context.Context, id id.User, updateFn domuser.UpdateFunc) error {
	args := m.Called(id, updateFn)
	return args.Error(0)
}

type mockSessionService struct {
	mock.Mock
}

func (m *mockSessionService) CreateForUser(_ context.Context, user *domuser.User) (*domsession.Session, error) {
	args := m.Called(user)
	return args.Get(0).(*domsession.Session), args.Error(1)
}

func (m *mockSessionService) DestroyForUser(_ context.Context, userID id.User) error {
	args := m.Called(userID)
	return args.Error(0)
}

func changePasswordMatchFunc(input, expected *domuser.User, expectedErr error) func(func(*domuser.User) (*domuser.User, error)) bool {
	return func(f func(*domuser.User) (*domuser.User, error)) bool {
		if input == nil {
			return true
		}
		resultUser, err := f(input)
		if !errors.Is(err, expectedErr) {
			return false
		}
		if expectedErr != nil {
			return true
		}
		return compareUsers(expected, resultUser)
	}
}

func compareUsers(expected, actual *domuser.User) bool {
	return expected.ReferrerID == actual.ReferrerID &&
		expected.Name == actual.Name &&
		expected.Email == actual.Email &&
		bytes.Equal(expected.PasswordHash, actual.PasswordHash) &&
		expected.Role == actual.Role &&
		expected.CreatedAt.Equal(actual.CreatedAt) &&
		expected.UpdatedAt.Equal(actual.UpdatedAt)
}

type mocks struct {
	hasher         *mockHasher
	timeSource     *mockTimeSource
	repository     *mockRepository
	sessionService *mockSessionService
}

func newMocks() mocks {
	return mocks{
		hasher:         &mockHasher{},
		timeSource:     &mockTimeSource{},
		repository:     &mockRepository{},
		sessionService: &mockSessionService{},
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.hasher.AssertExpectations(t)
	m.timeSource.AssertExpectations(t)
	m.repository.AssertExpectations(t)
	m.sessionService.AssertExpectations(t)
}

func newFactory(t *testing.T, hasher domuser.Hasher, timeSource timesource.TimeSource) domuser.Factory {
	t.Helper()
	userFactory, err := domuser.NewFactory(
		hasher,
		timeSource,
	)
	require.NoError(t, err)
	return userFactory
}

func newUser(t *testing.T, hasher domuser.Hasher, timeSource timesource.TimeSource) *domuser.User {
	t.Helper()
	factory := newFactory(t, hasher, timeSource)
	createUserInput := domuser.CreateUserInput{
		Name:       testUserName,
		Email:      testUserEmail,
		Password:   testUserPassword,
		ReferrerID: nil,
	}
	user, err := factory.NewUser(createUserInput, domuser.RoleUser)
	require.NoError(t, err)
	return user
}

func Test_NewService(t *testing.T) {
	mocks := newMocks()
	factory := newFactory(t, mocks.hasher, mocks.timeSource)
	expected := &Service{
		userFactory:    factory,
		userRepository: mocks.repository,
		sessionService: mocks.sessionService,
	}
	service, err := NewService(factory, mocks.repository, mocks.sessionService)
	assert.NoError(t, err)
	assert.Equal(t, expected, service)

	service, err = NewService(factory, mocks.repository, nil)
	assert.EqualError(t, err, "invalid session service")
	assert.Empty(t, service)

	service, err = NewService(factory, nil, nil)
	assert.EqualError(t, err, "invalid user repository")
	assert.Empty(t, service)
}

func Test_Service_Create(t *testing.T) {
	createUserInput := domuser.CreateUserInput{
		Name:       testUserName,
		Email:      testUserEmail,
		Password:   testUserPassword,
		ReferrerID: nil,
	}
	passwordHash := []byte("as46d4ad36a8d4a")
	now := time.Now()
	var user *domuser.User
	var session *domsession.Session

	type args struct {
		createUserInput domuser.CreateUserInput
	}
	tests := []struct {
		name            string
		args            args
		mocks           mocks
		expectedUser    *domuser.User
		expectedSession *domsession.Session
		expectedErr     error
	}{
		{
			name: "success",
			args: args{createUserInput: createUserInput},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				user = newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("Create", mock.MatchedBy(createUserMatchFunc(user))).Return(nil)
				sessionFactory, err := domsession.NewFactory(
					[]byte("abc123"),
					timesource.DefaultTimeSource{},
					time.Hour,
					time.Hour,
				)
				require.NoError(t, err)
				session, err = sessionFactory.NewSession(domsession.Claims{
					UserID: user.ID,
					Custom: domsession.CustomClaims{UserRole: user.Role},
				})
				require.NoError(t, err)
				mocks.sessionService.On("CreateForUser", mock.MatchedBy(createUserMatchFunc(user))).Return(session, nil)
				return mocks
			}(),
			expectedUser:    user,
			expectedSession: session,
			expectedErr:     nil,
		},
		{
			name: "failure:create-for-user",
			args: args{createUserInput: createUserInput},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				user := newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("Create", mock.MatchedBy(createUserMatchFunc(user))).Return(nil)
				mocks.sessionService.On("CreateForUser", mock.MatchedBy(createUserMatchFunc(user))).Return((*domsession.Session)(nil), errTest)
				return mocks
			}(),
			expectedUser:    nil,
			expectedSession: nil,
			expectedErr:     fmt.Errorf("creating session for user: %w", errTest),
		},
		{
			name: "failure:create",
			args: args{createUserInput: createUserInput},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				mocks.repository.On("Create", mock.MatchedBy(createUserMatchFunc(user))).Return(errTest)
				return mocks
			}(),
			expectedUser:    nil,
			expectedSession: nil,
			expectedErr:     fmt.Errorf("creating user: %w", errTest),
		},
		{
			name: "failure:create/referrer-not-found",
			args: args{createUserInput: createUserInput},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				mocks.repository.On("Create", mock.MatchedBy(createUserMatchFunc(user))).Return(domuser.ErrReferrerNotFound)
				return mocks
			}(),
			expectedUser:    nil,
			expectedSession: nil,
			expectedErr:     apierrors.NewBadRequestError(domuser.ErrReferrerNotFound, "creating user").WithPublicMessage(domuser.ErrReferrerNotFound.Error()),
		},
		{
			name: "failure:create/user-email-already-exists",
			args: args{createUserInput: createUserInput},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				mocks.repository.On("Create", mock.MatchedBy(createUserMatchFunc(user))).Return(domuser.ErrUserEmailAlreadyExists)
				return mocks
			}(),
			expectedUser:    nil,
			expectedSession: nil,
			expectedErr:     apierrors.NewAlreadyExistsError(domuser.ErrUserEmailAlreadyExists, "creating user").WithPublicMessage(domuser.ErrUserEmailAlreadyExists.Error()),
		},
		{
			name: "failure:new-user",
			args: func() args {
				return args{createUserInput: createUserInput}
			}(),
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return([]byte(nil), errTest)
				return mocks
			}(),
			expectedUser:    nil,
			expectedSession: nil,
			expectedErr:     fmt.Errorf("new user: %w", errTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.hasher, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			user, session, err := service.Create(context.Background(), tt.args.createUserInput)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedUser != nil {
				user.ID = tt.expectedUser.ID
			}
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedSession, session)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_Read(t *testing.T) {
	passwordHash := []byte("sa7daas24as24da4sa")
	var user *domuser.User

	type args struct {
		userID id.User
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domuser.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{userID: user.ID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				user = newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("Read", user.ID).Return(user, nil)
				return mocks
			}(),
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "failure:read",
			args: args{userID: user.ID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Read", user.ID).Return((*domuser.User)(nil), errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("reading user: %w", errTest),
		},
		{
			name: "failure:read/not-found",
			args: args{userID: user.ID},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Read", user.ID).Return((*domuser.User)(nil), domuser.ErrUserNotFound)
				return mocks
			}(),
			expected:    nil,
			expectedErr: apierrors.NewNotFoundError(domuser.ErrUserNotFound, "reading user").WithPublicMessage(domuser.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.hasher, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			user, err := service.Read(context.Background(), tt.args.userID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_ReadByEmail(t *testing.T) {
	var user *domuser.User

	type args struct {
		email types.Email
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domuser.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{email: testUserEmail},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return([]byte("sa7daas24as24da4sa"), nil)
				mocks.timeSource.On("Now").Return(time.Now())
				user = newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("ReadByEmail", user.Email).Return(user, nil)
				return mocks
			}(),
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "failure:read-by-email",
			args: args{email: testUserEmail},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("ReadByEmail", testUserEmail).Return((*domuser.User)(nil), errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("reading user by email: %w", errTest),
		},
		{
			name: "failure:read-by-email/not-found",
			args: args{email: testUserEmail},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("ReadByEmail", testUserEmail).Return((*domuser.User)(nil), domuser.ErrUserNotFound)
				return mocks
			}(),
			expected:    nil,
			expectedErr: apierrors.NewNotFoundError(domuser.ErrUserNotFound, "reading user by email").WithPublicMessage(domuser.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.hasher, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			user, err := service.ReadByEmail(context.Background(), tt.args.email)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_ReadByCredentials(t *testing.T) {
	passwordHash := []byte("sa7daas24as24da4sa")
	var user *domuser.User

	type args struct {
		email    types.Email
		password types.Password
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domuser.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				email:    testUserEmail,
				password: testUserPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				user = newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("ReadByEmail", user.Email).Return(user, nil)
				mocks.hasher.On("CompareHashAndPassword", passwordHash, []byte(testUserPassword)).Return(true)
				return mocks
			}(),
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "failure:match-password",
			args: args{
				email:    testUserEmail,
				password: testUserPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				user := newUser(t, mocks.hasher, mocks.timeSource)
				mocks.repository.On("ReadByEmail", user.Email).Return(user, nil)
				mocks.hasher.On("CompareHashAndPassword", passwordHash, []byte(testUserPassword)).Return(false)
				return mocks
			}(),
			expected:    nil,
			expectedErr: apierrors.NewUnauthorizedError(errors.New("invalid password"), "").WithPublicMessage("email or password is incorrect"),
		},
		{
			name: "failure:read-by-email",
			args: args{
				email:    testUserEmail,
				password: testUserPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("ReadByEmail", testUserEmail).Return((*domuser.User)(nil), errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("reading user by credentials: %w", errTest),
		},
		{
			name: "failure:read-by-email/not-found",
			args: args{
				email:    testUserEmail,
				password: testUserPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("ReadByEmail", testUserEmail).Return((*domuser.User)(nil), domuser.ErrUserNotFound)
				return mocks
			}(),
			expected:    nil,
			expectedErr: apierrors.NewUnauthorizedError(domuser.ErrUserNotFound, "reading user by credentials").WithPublicMessage("email or password is incorrect"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := newFactory(t, tt.mocks.hasher, tt.mocks.timeSource)
			service, err := NewService(factory, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			user, err := service.ReadByCredentials(context.Background(), tt.args.email, tt.args.password)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_ChangePassword(t *testing.T) {
	userID := id.NewUser()
	oldPassword := testUserPassword
	oldPasswordHash := []byte("sa7daas24as24da4sa")
	newPassword := types.Password("Random123")
	newPasswordHash := []byte("dfs45f4s6ds8sdsdf5")

	type args struct {
		userID      id.User
		oldPassword types.Password
		newPassword types.Password
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				userID:      userID,
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(oldPasswordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now()).Once()
				user := newUser(t, mocks.hasher, mocks.timeSource)
				updatedUser := *user
				updatedUser.PasswordHash = newPasswordHash
				updatedUser.UpdatedAt = user.UpdatedAt.Add(time.Hour)
				mocks.timeSource.On("Now").Return(updatedUser.UpdatedAt).Once()
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(true)
				mocks.hasher.On("HashPassword", []byte(newPassword)).Return(newPasswordHash, nil)
				mocks.repository.On("Update", userID, mock.MatchedBy(changePasswordMatchFunc(user, &updatedUser, nil))).Return(nil)
				mocks.sessionService.On("DestroyForUser", userID).Return(nil)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:destroy-sessions-for-user",
			args: args{
				userID:      userID,
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(oldPasswordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now()).Once()
				user := newUser(t, mocks.hasher, mocks.timeSource)
				updatedUser := *user
				updatedUser.PasswordHash = newPasswordHash
				updatedUser.UpdatedAt = user.UpdatedAt.Add(time.Hour)
				mocks.timeSource.On("Now").Return(updatedUser.UpdatedAt).Once()
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(true)
				mocks.hasher.On("HashPassword", []byte(newPassword)).Return(newPasswordHash, nil)
				mocks.repository.On("Update", userID, mock.MatchedBy(changePasswordMatchFunc(user, &updatedUser, nil))).Return(nil)
				mocks.sessionService.On("DestroyForUser", userID).Return(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("destroying sessions for user: %w", errTest),
		},
		{
			name: "failure:update",
			args: args{
				userID:      userID,
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Update", userID, mock.MatchedBy(changePasswordMatchFunc(nil, nil, nil))).Return(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("changing password: %w", errTest),
		},
		{
			name: "failure:update/invalid-user-password",
			args: args{
				userID:      userID,
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return(oldPasswordHash, nil)
				mocks.timeSource.On("Now").Return(time.Now())
				user := newUser(t, mocks.hasher, mocks.timeSource)
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(false)
				mocks.repository.On("Update", userID, mock.MatchedBy(changePasswordMatchFunc(user, nil, domuser.ErrInvalidUserPassword))).Return(domuser.ErrInvalidUserPassword)
				return mocks
			}(),
			expectedErr: apierrors.NewBadRequestError(domuser.ErrInvalidUserPassword, "changing password").WithPublicMessage(domuser.ErrInvalidUserPassword.Error()),
		},
		{
			name: "failure:update/user-not-found",
			args: args{
				userID: userID,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("Update", userID, mock.MatchedBy(changePasswordMatchFunc(nil, nil, nil))).Return(domuser.ErrUserNotFound)
				return mocks
			}(),
			expectedErr: apierrors.NewUnauthorizedError(domuser.ErrUserNotFound, "changing password").WithPublicMessage(domuser.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(domuser.Factory{}, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			err = service.ChangePassword(context.Background(), tt.args.userID, tt.args.oldPassword, tt.args.newPassword)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Service_List(t *testing.T) {
	var users []domuser.User

	tests := []struct {
		name        string
		mocks       mocks
		expected    []domuser.User
		expectedErr error
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(testUserPassword)).Return([]byte("4ds5fs"), nil)
				mocks.timeSource.On("Now").Return(time.Now())
				users = []domuser.User{
					*newUser(t, mocks.hasher, mocks.timeSource),
					*newUser(t, mocks.hasher, mocks.timeSource),
				}
				mocks.repository.On("List").Return(users, nil)
				return mocks
			}(),
			expected:    users,
			expectedErr: nil,
		},
		{
			name: "failure:list",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.repository.On("List").Return(([]domuser.User)(nil), errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("listing users: %w", errTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(domuser.Factory{}, tt.mocks.repository, tt.mocks.sessionService)
			assert.NoError(t, err)
			users, err := service.List(context.Background())
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, users)
			tt.mocks.assertExpectations(t)
		})
	}
}
