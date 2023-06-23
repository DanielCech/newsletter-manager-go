package newsletter

import (
	"errors"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"
	"time"
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

func (f Factory) NewNewsletterFromFields(
	id id.Newsletter,
	authorId id.Author,
	name string,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) *Newsletter {
	return &Newsletter{
		timeSource:  f.timeSource,
		ID:          id,
		AuthorID:    authorId,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func newFactoryValidate(timesource timesource.TimeSource) error {
	if timesource == nil {
		return errors.New("invalid time source")
	}
	return nil
}
