package v1

import (
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
)

func (h *Handler) AuthorSubscriptions(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) ([]model.Subscription, error) {
	// TODO
	return nil, nil
}

func (h *Handler) SubscribeToNewsletter(_ http.ResponseWriter, r *http.Request) ([]model.Subscription, error) {
	// TODO
	return nil, nil
}

func (h *Handler) UnsubscribeFromNewsletter(_ http.ResponseWriter, r *http.Request) ([]model.Subscription, error) {
	// TODO
	return nil, nil
}
