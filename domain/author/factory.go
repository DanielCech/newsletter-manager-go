package author

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

// Factory contains dependencies that are needed for authors creation.
type Factory struct {
	hasher     Hasher
	timeSource timesource.TimeSource
}

// NewFactory returns new instance of author Factory.
func NewFactory(timesource timesource.TimeSource) (Factory, error) {
	if err := newFactoryValidate(timesource); err != nil {
		return Factory{}, err
	}
	return Factory{
		timeSource: timesource,
	}, nil
}

// NewAuthorFromFields returns new instance of Author based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
func (f Factory) NewAuthorFromFields(
	id id.Author,
	name string,
	email string,
	passwordHash []byte,
	createdAt time.Time,
	updatedAt time.Time,
) *Author {
	return &Author{
		hasher:       f.hasher,
		timeSource:   f.timeSource,
		ID:           id,
		Name:         name,
		Email:        types.Email(email),
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func newFactoryValidate(timesource timesource.TimeSource) error {
	if timesource == nil {
		return errors.New("invalid time source")
	}
	return nil
}
