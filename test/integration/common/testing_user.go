package common

import (
	"context"
	"event-facematch-backend/test/integration/generate/swagger"
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

func (user *TestingUser) UpdateWith(resp swagger.SignInResp) {
	user.AuthorID = resp.User.Id
}