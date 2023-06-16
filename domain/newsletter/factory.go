package newsletter

import (
	"errors"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
)

// Factory contains dependencies that are needed for newsletters creation.
type Factory struct {
	timeSource timesource.TimeSource
}

// NewFactory returns new instance of newsletter Factory.
func NewFactory(timesource timesource.TimeSource) (Factory, error) {
	if err := newFactoryValidate(timesource); err != nil {
		return Factory{}, err
	}
	return Factory{
		timeSource: timesource,
	}, nil
}

// NewNewsletter returns new instance of Newsletter.
func (f Factory) NewNewsletter(createNewsletterInput CreateNewsletterInput) (*Newsletter, error) {
	now := f.timeSource.Now()
	newsletter := &Newsletter{
		timeSource:  f.timeSource,
		ID:          id.NewNewsletter(),
		AuthorID:    createNewsletterInput.AuthorID,
		Name:        createNewsletterInput.Name,
		Description: createNewsletterInput.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return newsletter, nil
}

// NewNewsletterFromFields returns new instance of Newsletter based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
// func (f Factory) NewNewsletterFromFields(
//	id id.Newsletter,
//	referrerID *id.Newsletter,
//	name string,
//	email string,
//	passwordHash []byte,
//	role string,
//	createdAt time.Time,
//	updatedAt time.Time,
// ) *Newsletter {
//	return &Newsletter{
//		hasher:       f.hasher,
//		timeSource:   f.timeSource,
//		ID:           id,
//		ReferrerID:   referrerID,
//		Name:         name,
//		Email:        types.Email(email),
//		PasswordHash: passwordHash,
//		Role:         Role(role),
//		CreatedAt:    createdAt,
//		UpdatedAt:    updatedAt,
//	}
// }

func newFactoryValidate(timesource timesource.TimeSource) error {
	if timesource == nil {
		return errors.New("invalid time source")
	}
	return nil
}
