package types

import (
	"errors"
	"strings"

	"strv-template-backend-go-api/types/validator"
)

// Email is a custom type for email address representation.
// Email addresses should be always printed in a lower case.
type Email string

// NewEmail returns new instance of Email.
func NewEmail(s string) (Email, error) {
	e := Email(strings.ToLower(s))
	if err := e.valid(); err != nil {
		return "", err
	}
	return e, nil
}

// UnmarshalText unmarshals text value into the Email.
func (e *Email) UnmarshalText(data []byte) error {
	*e = Email(strings.ToLower(string(data)))
	return e.valid()
}

// String returns email always in lower case.
func (e Email) String() string {
	return strings.ToLower(string(e))
}

// Valid checks whether email value is valid.
func (e Email) Valid() error {
	if string(e) != strings.ToLower(string(e)) {
		return errors.New("email is not lowercase")
	}
	return e.valid()
}

func (e Email) valid() error {
	if err := validator.Validate.Var(e, "email"); err != nil {
		return err
	}
	return nil
}
