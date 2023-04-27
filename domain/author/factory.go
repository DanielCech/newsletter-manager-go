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

// NewAuthor returns new instance of Author.
func (f Factory) NewAuthor(createAuthorInput CreateAuthorInput, role Role) (*Author, error) {
	passwordHash, err := f.hasher.HashPassword([]byte(createAuthorInput.Password))
	if err != nil {
		return nil, err
	}

	now := f.timeSource.Now()
	author := &Author{
		hasher:       f.hasher,
		timeSource:   f.timeSource,
		ID:           id.NewAuthor(),
		ReferrerID:   createAuthorInput.ReferrerID,
		Name:         createAuthorInput.Name,
		Email:        createAuthorInput.Email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = author.Valid(); err != nil {
		return nil, err
	}

	return author, nil
}

// NewAuthorFromFields returns new instance of Author based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
func (f Factory) NewAuthorFromFields(
	id id.Author,
	referrerID *id.Author,
	name string,
	email string,
	passwordHash []byte,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) *Author {
	return &Author{
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

func newFactoryValidate(timesource timesource.TimeSource) error {
	if timesource == nil {
		return errors.New("invalid time source")
	}
	return nil
}
