package v1

import (
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
)

func (h *Handler) ListNewsletterEmails(_ http.ResponseWriter, r *http.Request, input model.PathNewsletterInput) ([]model.FullEmail, error) {
	// TODO
	return nil, nil
}

func (h *Handler) CreateNewsletterEmail(_ http.ResponseWriter, r *http.Request, input model.PathNewsletterInput) (*model.Email, error) {
	// TODO
	return nil, nil
}
