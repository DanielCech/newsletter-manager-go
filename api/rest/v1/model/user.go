package model

import (
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"
)

// User consists of fields which describe an user.
type User struct {
	ID         id.User     `json:"id"`
	Name       string      `json:"name"`
	Email      types.Email `json:"email"`
	Role       string      `json:"role"`
	ReferrerID *id.User    `json:"referrerId,omitempty"`
}

// FromUser converts domain object to api object.
func FromUser(user *domuser.User) User {
	return User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}

// FromUsers converts domain object to api object.
func FromUsers(dusers []domuser.User) []User {
	users := make([]User, 0, len(dusers))
	for _, u := range dusers {
		users = append(users, User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  string(u.Role),
		})
	}
	return users
}

// CreateUserInput represents JSON body needed for creating a new user.
type CreateUserInput struct {
	Name       string         `json:"name" validate:"required"`
	Email      types.Email    `json:"email"`
	Password   types.Password `json:"password"`
	ReferrerID *id.User       `json:"referrerId"`
}

// CreateUserResp represents JSON response body of creating a new user.
type CreateUserResp struct {
	User    User    `json:"user"`
	Session Session `json:"session"`
}

// ChangeUserPasswordInput represents JSON body needed for changing the user password.
type ChangeUserPasswordInput struct {
	OldPassword types.Password `json:"oldPassword"`
	NewPassword types.Password `json:"newPassword"`
}
