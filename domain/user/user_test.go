package user

import (
	"testing"
	"time"

	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"

	"github.com/stretchr/testify/assert"
)

func Test_NewRole(t *testing.T) {
	role, err := NewRole("user")
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, role)

	role, err = NewRole("unknown")
	assert.Error(t, err)
	assert.Empty(t, role)
}

func Test_Role_Valid(t *testing.T) {
	assert.NoError(t, Role("user").Valid())
	assert.NoError(t, Role("admin").Valid())
	assert.Equal(t, ErrInvalidUserRole, Role("unknown").Valid())
}

func Test_Role_IsSufficientToRole(t *testing.T) {
	role := RoleAdmin
	assert.True(t, role.IsSufficientToRole(RoleAdmin))
	assert.True(t, role.IsSufficientToRole(RoleUser))

	role = RoleUser
	assert.True(t, role.IsSufficientToRole(RoleUser))
	assert.False(t, role.IsSufficientToRole(RoleAdmin))
}

func Test_User_Valid(t *testing.T) {
	user := User{
		ID:           id.NewUser(),
		Name:         "Jozko",
		Email:        "Dlouhy",
		PasswordHash: []byte("45sa6d5a63"),
		Role:         RoleUser,
	}
	assert.NoError(t, user.Valid())

	user.Role = "unknown"
	assert.Error(t, user.Valid())

	user.PasswordHash = nil
	assert.Equal(t, ErrInvalidUserPassword, user.Valid())

	user.Name = ""
	assert.Equal(t, ErrInvalidUserName, user.Valid())

	user.ID = id.User{}
	assert.Equal(t, ErrInvalidUserID, user.Valid())
}

func Test_User_MatchPassword(t *testing.T) {
	password := types.Password("Topsecret1")
	passwordHash := []byte("4asd47as4d7as545dsa54sa45d57asd4as")
	hasher := &mockHasher{}
	hasher.On("CompareHashAndPassword", passwordHash, []byte(password)).Return(true).Once()

	user := User{
		hasher:       hasher,
		PasswordHash: passwordHash,
	}
	assert.True(t, user.MatchPassword(password))

	hasher.On("CompareHashAndPassword", passwordHash, []byte(password)).Return(false).Once()
	assert.False(t, user.MatchPassword(password))

	hasher.AssertExpectations(t)
}

func Test_User_ChangePassword(t *testing.T) {
	oldPassword := types.Password("Topsecret1")
	oldPasswordHash := []byte("4dsa65d4a5dsa45sda45as3")
	newPassword := types.Password("Topsecret2")
	newPasswordHash := []byte("4asd47as4d7as545dsa54sa45d57asd4as")
	now := time.Now()

	type args struct {
		oldPassword types.Password
		newPassword types.Password
	}
	tests := []struct {
		name     string
		args     args
		mocks    mocks
		assertFn func(*testing.T, User, args)
	}{
		{
			name: "success",
			args: args{
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(true)
				mocks.hasher.On("HashPassword", []byte(newPassword)).Return(newPasswordHash, nil)
				mocks.timeSource.On("Now").Return(now)
				return mocks
			}(),
			assertFn: func(t *testing.T, user User, args args) {
				t.Helper()
				err := user.ChangePassword(args.oldPassword, args.newPassword)
				assert.NoError(t, err)
				assert.Equal(t, newPasswordHash, user.PasswordHash)
				assert.Equal(t, now, user.UpdatedAt)
			},
		},
		{
			name: "failure:hash-password",
			args: args{
				oldPassword: oldPassword,
				newPassword: newPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(true)
				mocks.hasher.On("HashPassword", []byte(newPassword)).Return([]byte(nil), errTest)
				return mocks
			}(),
			assertFn: func(t *testing.T, user User, args args) {
				t.Helper()
				err := user.ChangePassword(args.oldPassword, args.newPassword)
				assert.Equal(t, errTest, err)
				assert.Equal(t, oldPasswordHash, user.PasswordHash)
			},
		},
		{
			name: "failure:match-password",
			args: args{
				oldPassword: oldPassword,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.hasher.On("CompareHashAndPassword", oldPasswordHash, []byte(oldPassword)).Return(false)
				return mocks
			}(),
			assertFn: func(t *testing.T, user User, args args) {
				t.Helper()
				err := user.ChangePassword(args.oldPassword, args.newPassword)
				assert.Equal(t, ErrInvalidUserPassword, err)
				assert.Equal(t, oldPasswordHash, user.PasswordHash)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{
				hasher:       tt.mocks.hasher,
				timeSource:   tt.mocks.timeSource,
				PasswordHash: oldPasswordHash,
			}
			tt.assertFn(t, user, tt.args)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_NewCreateUserInput(t *testing.T) {
	name := "Jozko"
	email := types.Email("jozko.dlouhy@gmail.com")
	password := types.Password("Topsecret1")
	expected := CreateUserInput{
		Name:       name,
		Email:      email,
		Password:   password,
		ReferrerID: nil,
	}

	createUserInput, err := NewCreateUserInput(name, email, password, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, createUserInput)

	createUserInput, err = NewCreateUserInput(name, email, "invalid", nil)
	assert.Error(t, err)
	assert.Empty(t, createUserInput)
}

func Test_CreateUserInput_Valid(t *testing.T) {
	createUserInput := CreateUserInput{
		Name:     "Jozko",
		Email:    types.Email("jozko.dlouhy@gmail.com"),
		Password: types.Password("Topsecret1"),
	}
	err := createUserInput.Valid()
	assert.NoError(t, err)

	createUserInput.Password = "invalid"
	err = createUserInput.Valid()
	assert.Error(t, err)

	createUserInput.Email = "invalid"
	err = createUserInput.Valid()
	assert.EqualError(t, err, "Key: '' Error:Field validation for '' failed on the 'email' tag")

	createUserInput.Name = ""
	err = createUserInput.Valid()
	assert.Equal(t, ErrInvalidUserName, err)
}
