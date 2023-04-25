package model

import (
	"fmt"
	"io"

	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	customvalidator "strv-template-backend-go-api/types/validator"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

const passwordResponse = "*******"

// MarshalPassword marshals types.Password to string.
func MarshalPassword(_ types.Password) graphql.Marshaler {
	return graphql.MarshalString(passwordResponse)
}

// UnmarshalPassword unmarshals value into the types.Password.
func UnmarshalPassword(v any) (types.Password, error) {
	switch v := v.(type) {
	case string:
		var password types.Password
		if err := password.UnmarshalText([]byte(v)); err != nil {
			const errMsg = "password validation failed"
			return "", apierrors.NewInvalidBodyError(err, errMsg).WithPublicMessage(errMsg)
		}
		return password, nil
	default:
		const publicErrMsg = "password validation failed"
		err := fmt.Errorf("%T must be a string", v)
		return "", apierrors.NewInvalidBodyError(err, publicErrMsg).WithPublicMessage(publicErrMsg)
	}
}

// Email is a custom scalar for email representation.
type Email string

// MarshalGQL marshals Email to string.
func (e Email) MarshalGQL(w io.Writer) {
	graphql.MarshalString(string(e)).MarshalGQL(w)
}

// UnmarshalGQL unmarshals value into the Email.
func (e *Email) UnmarshalGQL(v any) error {
	s, ok := v.(string)
	if !ok {
		const publicErrMsg = "email validation failed"
		err := fmt.Errorf("%T must be a string", v)
		return apierrors.NewInvalidBodyError(err, publicErrMsg).WithPublicMessage(publicErrMsg)
	}

	if err := customvalidator.Validate.Var(s, "email"); err != nil {
		const errMsg = "email validation failed"
		return apierrors.NewInvalidBodyError(err, errMsg).WithPublicMessage(errMsg)
	}

	*e = Email(s)
	return nil
}

// User contains all fields which describes an user.
type User struct {
	ID         uuid.UUID  `json:"id"`
	ReferrerID *uuid.UUID `json:"-"`
	Name       string     `json:"name"`
	Email      Email      `json:"email"`
	Role       Role       `json:"role"`
}

// FromUser converts domain object to graphql object.
func FromUser(user *domuser.User) *User {
	return &User{
		ID:         uuid.UUID(user.ID),
		ReferrerID: (*uuid.UUID)(user.ReferrerID),
		Name:       user.Name,
		Email:      Email(user.Email.String()),
		Role:       Role(user.Role),
	}
}

// FromUsers converts domain object to graphql object.
func FromUsers(domUsers []domuser.User) []User {
	users := make([]User, 0, len(domUsers))
	for i := range domUsers {
		users = append(users, *FromUser(&domUsers[i]))
	}
	return users
}

// FromSession converts domain object to graphql object.
func FromSession(session *domsession.Session) *RefreshSessionResponse {
	return &RefreshSessionResponse{
		AccessToken:           session.AccessToken.SignedData,
		AccessTokenExpiresAt:  session.AccessToken.ExpiresAt,
		RefreshToken:          string(session.RefreshToken.ID),
		RefreshTokenExpiresAt: session.RefreshToken.ExpiresAt,
	}
}
