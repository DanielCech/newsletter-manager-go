package user

import (
	"errors"
	"time"

	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
)

// Hasher describes object which is capable of password hashing and comparing.
type Hasher interface {
	HashPassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hash, password []byte) bool
}

// Factory contains dependencies that are needed for users creation.
type Factory struct {
	hasher     Hasher
	timeSource timesource.TimeSource
}

// NewFactory returns new instance of user Factory.
func NewFactory(hasher Hasher, timesource timesource.TimeSource) (Factory, error) {
	if err := newFactoryValidate(hasher, timesource); err != nil {
		return Factory{}, err
	}
	return Factory{
		hasher:     hasher,
		timeSource: timesource,
	}, nil
}

// NewUser returns new instance of User.
func (f Factory) NewUser(createUserInput CreateUserInput, role Role) (*User, error) {
	passwordHash, err := f.hasher.HashPassword([]byte(createUserInput.Password))
	if err != nil {
		return nil, err
	}

	now := f.timeSource.Now()
	user := &User{
		hasher:       f.hasher,
		timeSource:   f.timeSource,
		ID:           id.NewUser(),
		ReferrerID:   createUserInput.ReferrerID,
		Name:         createUserInput.Name,
		Email:        createUserInput.Email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = user.Valid(); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserFromFields returns new instance of User based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
func (f Factory) NewUserFromFields(
	id id.User,
	referrerID *id.User,
	name string,
	email string,
	passwordHash []byte,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		hasher:       f.hasher,
		timeSource:   f.timeSource,
		ID:           id,
		ReferrerID:   referrerID,
		Name:         name,
		Email:        types.Email(email),
		PasswordHash: passwordHash,
		Role:         Role(role),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func newFactoryValidate(hasher Hasher, timesource timesource.TimeSource) error {
	if hasher == nil {
		return errors.New("invalid hasher")
	}
	if timesource == nil {
		return errors.New("invalid time source")
	}
	return nil
}
