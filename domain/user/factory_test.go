package user

import (
	"errors"
	"testing"
	"time"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

func ptr[T any](v T) *T {
	return &v
}

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

type mocks struct {
	hasher     *mockHasher
	timeSource *mockTimeSource
}

func newMocks() mocks {
	return mocks{
		hasher:     &mockHasher{},
		timeSource: &mockTimeSource{},
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.hasher.AssertExpectations(t)
	m.timeSource.AssertExpectations(t)
}

func Test_NewFactory(t *testing.T) {
	mocks := newMocks()
	expected := Factory{
		hasher:     mocks.hasher,
		timeSource: mocks.timeSource,
	}
	factory, err := NewFactory(mocks.hasher, mocks.timeSource)
	assert.NoError(t, err)
	assert.Equal(t, expected, factory)

	factory, err = NewFactory(mocks.hasher, nil)
	assert.Error(t, err)
	assert.Empty(t, factory)
}

func Test_Factory_NewUser(t *testing.T) {
	createUserInput := CreateUserInput{
		Name:       "Jozko",
		Email:      types.Email("jozko.dlouhy@gmail.com"),
		Password:   types.Password("Topsecret1"),
		ReferrerID: ptr(id.NewUser()),
	}
	passwordHash := []byte("sa546d4a6sdas63da58ds6a38d4a")
	role := RoleUser
	now := time.Now()

	type args struct {
		createUserInput CreateUserInput
		role            Role
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		assertFn    func(*testing.T, *User)
		assertErrFn func(*testing.T, error)
	}{
		{
			name: "success",
			args: args{
				createUserInput: createUserInput,
				role:            role,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				return mocks
			}(),
			assertFn: func(t *testing.T, user *User) {
				t.Helper()
				assert.NotNil(t, user.hasher)
				assert.NotNil(t, user.timeSource)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, createUserInput.ReferrerID, user.ReferrerID)
				assert.Equal(t, createUserInput.Name, user.Name)
				assert.Equal(t, createUserInput.Email, user.Email)
				assert.Equal(t, passwordHash, user.PasswordHash)
				assert.Equal(t, role, user.Role)
				assert.Equal(t, now, user.CreatedAt)
				assert.Equal(t, now, user.UpdatedAt)
			},
			assertErrFn: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
		},
		{
			name: "failure:user-valid",
			args: args{
				createUserInput: CreateUserInput{
					Password: createUserInput.Password,
				},
				role: role,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return(passwordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				return mocks
			}(),
			assertFn: func(t *testing.T, user *User) {
				t.Helper()
				assert.Empty(t, user)
			},
			assertErrFn: func(t *testing.T, err error) {
				t.Helper()
				assert.Error(t, err)
			},
		},
		{
			name: "failure:hash-password",
			args: args{
				createUserInput: createUserInput,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("HashPassword", []byte(createUserInput.Password)).Return([]byte(nil), errTest)
				return mocks
			}(),
			assertFn: func(t *testing.T, user *User) {
				t.Helper()
				assert.Empty(t, user)
			},
			assertErrFn: func(t *testing.T, err error) {
				t.Helper()
				assert.Equal(t, errTest, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := NewFactory(tt.mocks.hasher, tt.mocks.timeSource)
			require.NoError(t, err)
			user, err := factory.NewUser(tt.args.createUserInput, tt.args.role)
			tt.assertErrFn(t, err)
			tt.assertFn(t, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Factory_NewUserFromFields(t *testing.T) {
	userID := id.NewUser()
	referrerID := id.NewUser()
	name := "Jozko"
	email := "jozko.dlouhy@gmail.com"
	passwordHash := []byte("sa546d4a6sdas63da58ds6a38d4a")
	role := "user"
	now := time.Now()
	mocks := newMocks()
	factory, err := NewFactory(mocks.hasher, mocks.timeSource)
	require.NoError(t, err)
	expected := &User{
		hasher:       mocks.hasher,
		timeSource:   mocks.timeSource,
		ID:           userID,
		ReferrerID:   &referrerID,
		Name:         name,
		Email:        types.Email(email),
		PasswordHash: passwordHash,
		Role:         Role(role),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	user := factory.NewUserFromFields(userID, &referrerID, name, email, passwordHash, role, now, now)
	assert.Equal(t, expected, user)
}

func Test_newFactoryValidate(t *testing.T) {
	mocks := newMocks()
	err := newFactoryValidate(mocks.hasher, mocks.timeSource)
	assert.NoError(t, err)

	err = newFactoryValidate(mocks.hasher, nil)
	assert.EqualError(t, err, "invalid time source")

	err = newFactoryValidate(nil, nil)
	assert.EqualError(t, err, "invalid hasher")
}
