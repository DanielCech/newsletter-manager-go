package v1

import (
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
	domsubscription "newsletter-manager-go/domain/subscription"
	apierrors "newsletter-manager-go/types/errors"
)

// CreateSubscription creates new subscription.
func (h *Handler) CreateSubscription(_ http.ResponseWriter, r *http.Request, input model.CreateSubscriptionInput) (*model.CreateSubscriptionResp, error) {
	createSubscriptionInput, err := domsubscription.NewCreateSubscriptionInput(
		input.Name,
		input.Email,
		input.Password,
		input.ReferrerID,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create subscription input").WithPublicMessage(err.Error())
	}
	subscription, session, err := h.newsletterService.Create(r.Context(), createSubscriptionInput)
	if err != nil {
		return nil, fmt.Errorf("creating subscription: %w", err)
	}
	createSubscriptionResp := model.CreateSubscriptionResp{
		Subscription: model.FromSubscription(subscription),
		Session:      model.FromSession(session),
	}
	return &createSubscriptionResp, nil
}

// ListSubscriptions returns all existing subscriptions.
func (h *Handler) ListSubscriptions(_ http.ResponseWriter, r *http.Request) ([]model.Subscription, error) {
	subscriptions, err := h.newsletterService.List(r.Context())
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}
	return model.FromSubscriptions(subscriptions), nil
}
