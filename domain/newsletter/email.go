package newsletter

import (
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Email consists of fields which describe an email.
type Email struct {
	ID    id.Email    `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}

// CreateNewsletterInput consists of fields required for creation of a new newsletter.
type CreateEmailInput struct {
	Name       string
	Email      types.Email
	Password   types.Password
	ReferrerID *id.Newsletter
}

// NewCreateNewsletterInput returns new instance of CreateNewsletterInput.
func NewCreateEmailInput(
	name string,
	email types.Email,
	password types.Password,
) (CreateNewsletterInput, error) {
	createNewsletterInput := CreateNewsletterInput{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := createNewsletterInput.Valid(); err != nil {
		return CreateNewsletterInput{}, err
	}
	return createNewsletterInput, nil
}

// Valid checks whether input fields are valid.
func (u CreateEmailInput) Valid() error {
	if len(u.Name) == 0 {
		return ErrInvalidNewsletterName
	}
	if err := u.Email.Valid(); err != nil {
		return err
	}
	return u.Password.Valid()
}
