package author

import (
	"errors"
	"time"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
)

var (
	ErrInvalidAuthorID       = errors.New("invalid author id")
	ErrInvalidAuthorName     = errors.New("invalid author name")
	ErrInvalidAuthorPassword = errors.New("invalid author password")
	ErrInvalidAuthorRole     = errors.New("invalid author role")
)

// Author consists of fields which describe a author.
type Author struct {
	hasher     Hasher
	timeSource timesource.TimeSource

	ID                     id.Author
	Name                   string
	Email                  types.Email
	PasswordHash           []byte
	AccessToken            string
	AccessTokenExpiration  string
	RefreshToken           string
	RefreshTokenExpiration string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (u *Author) Valid() error {
	if u.ID.Empty() {
		return ErrInvalidAuthorID
	}
	if len(u.Name) == 0 {
		return ErrInvalidAuthorName
	}
	if len(u.PasswordHash) == 0 {
		return ErrInvalidAuthorPassword
	}
	return nil
}

// MatchPassword compares author password hash with the given password.
func (u *Author) MatchPassword(password types.Password) bool {
	return u.hasher.CompareHashAndPassword(u.PasswordHash, []byte(password))
}

// ChangePassword checks whether the author password hash corresponds with the old password.
// If it does, author password hash is updated based on the new password.
func (u *Author) ChangePassword(oldPassword, newPassword types.Password) error {
	if !u.MatchPassword(oldPassword) {
		return ErrInvalidAuthorPassword
	}
	newPasswordHash, err := u.hasher.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = u.timeSource.Now()
	return nil
}

// CreateAuthorInput consists of fields required for creation of a new author.
type CreateAuthorInput struct {
	Name     string
	Email    types.Email
	Password types.Password
}

// NewCreateAuthorInput returns new instance of CreateAuthorInput.
func NewCreateAuthorInput(
	name string,
	email types.Email,
	password types.Password,
) (CreateAuthorInput, error) {
	createAuthorInput := CreateAuthorInput{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := createAuthorInput.Valid(); err != nil {
		return CreateAuthorInput{}, err
	}
	return createAuthorInput, nil
}

// Valid checks whether input fields are valid.
func (u CreateAuthorInput) Valid() error {
	if len(u.Name) == 0 {
		return ErrInvalidAuthorName
	}
	if err := u.Email.Valid(); err != nil {
		return err
	}
	return u.Password.Valid()
}

type GetNewsletterInput struct {
	AuthorID   id.Author
	Newsletter id.Newsletter
}

type NewsletterIDInput struct {
	NewsletterID id.Newsletter `json:"newsletterId"`
}
