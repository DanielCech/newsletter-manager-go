package newsletter

import (
	"errors"
	"time"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
)

var (
	ErrInvalidNewsletterID       = errors.New("invalid newsletter id")
	ErrInvalidNewsletterName     = errors.New("invalid newsletter name")
	ErrInvalidNewsletterPassword = errors.New("invalid newsletter password")
	ErrInvalidNewsletterRole     = errors.New("invalid newsletter role")
)

const (
	RoleNewsletter Role = "newsletter"
	RoleAdmin      Role = "admin"
)

// Role represents newsletter role.
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
//   - RoleNewsletter
//   - RoleAdmin
func (u Role) Valid() error {
	switch u {
	case RoleAdmin, RoleNewsletter:
		return nil
	}
	return ErrInvalidNewsletterRole
}

// IsSufficientToRole checks whether role is sufficient to the given one.
func (u Role) IsSufficientToRole(role Role) bool {
	switch role {
	case RoleAdmin:
		if u == RoleAdmin {
			return true
		}
	case RoleNewsletter:
		if u == RoleAdmin || u == RoleNewsletter {
			return true
		}
	}
	return false
}

// Newsletter consists of fields which describe a newsletter.
type Newsletter struct {
	hasher     Hasher
	timeSource timesource.TimeSource

	ID           id.Newsletter
	ReferrerID   *id.Newsletter
	Name         string
	Email        types.Email
	PasswordHash []byte
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Valid checks whether newsletter fields are valid.
func (u *Newsletter) Valid() error {
	if u.ID.Empty() {
		return ErrInvalidNewsletterID
	}
	if len(u.Name) == 0 {
		return ErrInvalidNewsletterName
	}
	if len(u.PasswordHash) == 0 {
		return ErrInvalidNewsletterPassword
	}
	return u.Role.Valid()
}

// MatchPassword compares newsletter password hash with the given password.
func (u *Newsletter) MatchPassword(password types.Password) bool {
	return u.hasher.CompareHashAndPassword(u.PasswordHash, []byte(password))
}

// ChangePassword checks whether the newsletter password hash corresponds with the old password.
// If it does, newsletter password hash is updated based on the new password.
func (u *Newsletter) ChangePassword(oldPassword, newPassword types.Password) error {
	if !u.MatchPassword(oldPassword) {
		return ErrInvalidNewsletterPassword
	}
	newPasswordHash, err := u.hasher.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = u.timeSource.Now()
	return nil
}

// CreateNewsletterInput consists of fields required for creation of a new newsletter.
type CreateNewsletterInput struct {
	Name       string
	Email      types.Email
	Password   types.Password
	ReferrerID *id.Newsletter
}

// NewCreateNewsletterInput returns new instance of CreateNewsletterInput.
func NewCreateNewsletterInput(
	name string,
	email types.Email,
	password types.Password,
	referrerID *id.Newsletter,
) (CreateNewsletterInput, error) {
	createNewsletterInput := CreateNewsletterInput{
		Name:       name,
		Email:      email,
		Password:   password,
		ReferrerID: referrerID,
	}
	if err := createNewsletterInput.Valid(); err != nil {
		return CreateNewsletterInput{}, err
	}
	return createNewsletterInput, nil
}

// Valid checks whether input fields are valid.
func (u CreateNewsletterInput) Valid() error {
	if len(u.Name) == 0 {
		return ErrInvalidNewsletterName
	}
	if err := u.Email.Valid(); err != nil {
		return err
	}
	return u.Password.Valid()
}
