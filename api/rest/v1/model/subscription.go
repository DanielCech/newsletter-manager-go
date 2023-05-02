package model

import (
	"newsletter-manager-go/types"
)

type Subscription struct {
	Email        string `json:"email"`
	NewsletterID int    `json:"newsletterId"`
}

type FullSubscription struct {
	Email        string `json:"email"`
	NewsletterID int    `json:"newsletterId"`
	Token        string `json:"token"`
}

//// FromSubscription converts domain object to api object.
//func FromDomainSubscription(subscription *domnewsletter.Subscription) Subscription {
//	return Subscription{
//		ID:    subscription.ID,
//		Name:  subscription.Name,
//		Email: subscription.Email,
//	}
//}
//
//// FromSubscriptions converts domain object to api object.
//func FromDomainSubscriptions(dsubscriptions []domnewsletter.Subscription) []Subscription {
//	subscriptions := make([]Subscription, 0, len(dsubscriptions))
//	for _, u := range dsubscriptions {
//		subscriptions = append(subscriptions, Subscription{
//			ID:    u.ID,
//			Name:  u.Name,
//			Email: u.Email,
//		})
//	}
//	return subscriptions
//}

// CreateSubscriptionInput represents JSON body needed for creating a new subscription.
type CreateSubscriptionInput struct {
	Name     string         `json:"name" validate:"required"`
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}
