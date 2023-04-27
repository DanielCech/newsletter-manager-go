package model

import (
	domsubscription "newsletter-manager-go/domain/subscription"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Subscription consists of fields which describe an subscription.
type Subscription struct {
	ID    id.Subscription `json:"id"`
	Name  string          `json:"name"`
	Email types.Email     `json:"email"`
}

// FromSubscription converts domain object to api object.
func FromSubscription(subscription *domsubscription.Subscription) Subscription {
	return Subscription{
		ID:    subscription.ID,
		Name:  subscription.Name,
		Email: subscription.Email,
	}
}

// FromSubscriptions converts domain object to api object.
func FromSubscriptions(dsubscriptions []domsubscription.Subscription) []Subscription {
	subscriptions := make([]Subscription, 0, len(dsubscriptions))
	for _, u := range dsubscriptions {
		subscriptions = append(subscriptions, Subscription{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return subscriptions
}

// CreateSubscriptionInput represents JSON body needed for creating a new subscription.
type CreateSubscriptionInput struct {
	Name       string           `json:"name" validate:"required"`
	Email      types.Email      `json:"email"`
	Password   types.Password   `json:"password"`
	ReferrerID *id.Subscription `json:"referrerId"`
}

// CreateSubscriptionResp represents JSON response body of creating a new subscription.
type CreateSubscriptionResp struct {
	Subscription Subscription `json:"subscription"`
	Session      Session      `json:"session"`
}
