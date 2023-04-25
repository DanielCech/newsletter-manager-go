package user

import (
	"errors"
	"time"

	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"
	"strv-template-backend-go-api/util/timesource"
)

var (
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidUserName     = errors.New("invalid user name")
	ErrInvalidUserPassword = errors.New("invalid user password")
	ErrInvalidUserRole     = errors.New("invalid user role")
)

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Role represents user role.
type Role string

// NewRole returns new instance of Role.
func NewRole(r string) (Role, error) {
	role := Role(r)
	if err := role.Valid(); err != nil {
		return "", err
	}
	return role, nil
}

// Valid checks whether role is valid.
// Possible values are:
//   - RoleUser
//   - RoleAdmin
func (u Role) Valid() error {
	switch u {
	case RoleAdmin, RoleUser:
		return nil
	}
	return ErrInvalidUserRole
}

// IsSufficientToRole checks whether role is sufficient to the given one.
func (u Role) IsSufficientToRole(role Role) bool {
	switch role {
	case RoleAdmin:
		if u == RoleAdmin {
			return true
		}
	case RoleUser:
		if u == RoleAdmin || u == RoleUser {
			return true
		}
	}
	return false
}

// User consists of fields which describe a user.
type User struct {
	hasher     Hasher
	timeSource timesource.TimeSource

	ID           id.User
	ReferrerID   *id.User
	Name         string
	Email        types.Email
	PasswordHash []byte
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Valid checks whether user fields are valid.
func (u *User) Valid() error {
	if u.ID.Empty() {
		return ErrInvalidUserID
	}
	if len(u.Name) == 0 {
		return ErrInvalidUserName
	}
	if len(u.PasswordHash) == 0 {
		return ErrInvalidUserPassword
	}
	return u.Role.Valid()
}

// MatchPassword compares user password hash with the given password.
func (u *User) MatchPassword(password types.Password) bool {
	return u.hasher.CompareHashAndPassword(u.PasswordHash, []byte(password))
}

// ChangePassword checks whether the user password hash corresponds with the old password.
// If it does, user password hash is updated based on the new password.
func (u *User) ChangePassword(oldPassword, newPassword types.Password) error {
	if !u.MatchPassword(oldPassword) {
		return ErrInvalidUserPassword
	}
	newPasswordHash, err := u.hasher.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = u.timeSource.Now()
	return nil
}

// CreateUserInput consists of fields required for creation of a new user.
type CreateUserInput struct {
	Name       string
	Email      types.Email
	Password   types.Password
	ReferrerID *id.User
}

// NewCreateUserInput returns new instance of CreateUserInput.
func NewCreateUserInput(
	name string,
	email types.Email,
	password types.Password,
	referrerID *id.User,
) (CreateUserInput, error) {
	createUserInput := CreateUserInput{
		Name:       name,
		Email:      email,
		Password:   password,
		ReferrerID: referrerID,
	}
	if err := createUserInput.Valid(); err != nil {
		return CreateUserInput{}, err
	}
	return createUserInput, nil
}

// Valid checks whether input fields are valid.
func (u CreateUserInput) Valid() error {
	if len(u.Name) == 0 {
		return ErrInvalidUserName
	}
	if err := u.Email.Valid(); err != nil {
		return err
	}
	return u.Password.Valid()
}
