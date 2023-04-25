package types

import (
	"errors"
	"unicode"
)

const (
	sensitiveOutput = "********"

	minPassLen = 8
	maxPassLen = 255
)

var (
	ErrMissingUpperCaseCharacter = errors.New("missing upper case character")
	ErrMissingLowerCaseCharacter = errors.New("missing lower case character")
	ErrMissingDigitCharacter     = errors.New("missing digit character")
	ErrInvalidPasswordLength     = errors.New("invalid password length")
)

// Password is a custom type for password representation.
// Password defines exact rules which have to be matched:
//   - Upper case character.
//   - Lower case character.
//   - Digit character.
//   - Minimum length of 8 characters.
//   - Maximum length of 255 characters.
type Password string

// NewPassword returns new instance of Password.
func NewPassword(p string) (Password, error) {
	password := Password(p)
	if err := password.Valid(); err != nil {
		return "", err
	}
	return password, nil
}

// MarshalText marshals Password to []byte.
func (Password) MarshalText() ([]byte, error) {
	return []byte(sensitiveOutput), nil
}

// UnmarshalText unmarshals text value into the Password.
func (p *Password) UnmarshalText(text []byte) error {
	*p = Password(text)
	return p.Valid()
}

// String returns masked Password by * characters.
// Password has to be never printed in the original form.
func (Password) String() string {
	return sensitiveOutput
}

// Valid checks exact rules which have to be matched:
//   - Upper case character.
//   - Lower case character.
//   - Digit character.
//   - Minimum length of 8 characters.
//   - Maximum length of 255 characters.
func (p Password) Valid() error {
	var err error
	if len(p) < minPassLen || len(p) > maxPassLen {
		err = errors.Join(err, ErrInvalidPasswordLength)
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range p {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasUpper {
		err = errors.Join(err, ErrMissingUpperCaseCharacter)
	}
	if !hasLower {
		err = errors.Join(err, ErrMissingLowerCaseCharacter)
	}
	if !hasDigit {
		err = errors.Join(err, ErrMissingDigitCharacter)
	}
	return err
}
