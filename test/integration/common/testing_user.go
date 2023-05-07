package common

import (
	"context"
	"newsletter-manager-go/test/integration/generate/swagger"

	"github.com/google/uuid"
)

type TestingUser struct {
	FirebaseID uuid.UUID
	AuthorID   string
	Context    context.Context
}

func NewUser() TestingUser {
	firebaseID := uuid.New()

	return TestingUser{
		FirebaseID: firebaseID,
		Context:    context.WithValue(context.Background(), swagger.ContextAccessToken, firebaseID.String()),
	}
}
