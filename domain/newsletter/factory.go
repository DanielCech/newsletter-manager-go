package newsletter

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

// Factory contains dependencies that are needed for newsletters creation.
type Factory struct {
	hasher     Hasher
	timeSource timesource.TimeSource
}

// NewFactory returns new instance of newsletter Factory.
func NewFactory(hasher Hasher, timesource timesource.TimeSource) (Factory, error) {
	if err := newFactoryValidate(hasher, timesource); err != nil {
		return Factory{}, err
	}
	return Factory{
		hasher:     hasher,
		timeSource: timesource,
	}, nil
}

// NewNewsletter returns new instance of Newsletter.
func (f Factory) NewNewsletter(createNewsletterInput CreateNewsletterInput, role Role) (*Newsletter, error) {
	passwordHash, err := f.hasher.HashPassword([]byte(createNewsletterInput.Password))
	if err != nil {
		return nil, err
	}

	now := f.timeSource.Now()
	newsletter := &Newsletter{
		hasher:       f.hasher,
		timeSource:   f.timeSource,
		ID:           id.NewNewsletter(),
		ReferrerID:   createNewsletterInput.ReferrerID,
		Name:         createNewsletterInput.Name,
		Email:        createNewsletterInput.Email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = newsletter.Valid(); err != nil {
		return nil, err
	}

	return newsletter, nil
}

// NewNewsletterFromFields returns new instance of Newsletter based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
func (f Factory) NewNewsletterFromFields(
	id id.Newsletter,
	referrerID *id.Newsletter,
	name string,
	email string,
	passwordHash []byte,
	role string,
	createdAt time.Time,
	updatedAt time.Time,
) *Newsletter {
	return &Newsletter{
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
