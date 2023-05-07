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

// CreateSubscriptionInput represents JSON body needed for creating a new subscription.
type CreateSubscriptionInput struct {
	Name     string         `json:"name" validate:"required"`
	Email    types.Email    `json:"email"`
	Password types.Password `json:"password"`
}
