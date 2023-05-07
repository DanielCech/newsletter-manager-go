package v1

import (
	"net/http"
	"newsletter-manager-go/api/rest/v1/model"
)

func (h *Handler) CreateNewsletter(_ http.ResponseWriter, r *http.Request, input model.CreateNewsletterInput) (*model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) ListNewsletters(_ http.ResponseWriter, r *http.Request) ([]model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) GetAuthorNewsletters(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) ([]model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) GetNewsletter(_ http.ResponseWriter, r *http.Request, input model.GetNewsletterInput) (*model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) UpdateNewsletter(_ http.ResponseWriter, r *http.Request, input model.NewsletterIDInput) (*model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) DeleteNewsletter(_ http.ResponseWriter, r *http.Request, input model.NewsletterIDInput) error {
	// TODO
	return nil
}
