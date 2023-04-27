package auth

import (
	"context"
	"event-facematch-backend/firebase"
	_ "event-facematch-backend/types"
	id "event-facematch-backend/types/id"
)

type IntegrationMockTokenParser struct {
	userIDMappings map[string]id.User
}

func NewIntegrationMockTokenParser() IntegrationMockTokenParser {
	return IntegrationMockTokenParser{userIDMappings: make(map[string]id.User)}
}

func (p *IntegrationMockTokenParser) VerifyIDToken(ctx context.Context, idToken string) (firebase.VerifiedToken, error) {
	userID, ok := p.userIDMappings[idToken]
	if !ok {
		return firebase.VerifiedToken{
			UserID:         nil,
			FirebaseUserID: idToken,
		}, nil
	}

	return firebase.VerifiedToken{
		UserID:         &userID,
		FirebaseUserID: idToken,
	}, nil
}

func (p *IntegrationMockTokenParser) SetUserIDCustomClaim(ctx context.Context, firebaseUserID string, userID id.User) error {
	p.userIDMappings[firebaseUserID] = userID
	return nil
}

func (p *IntegrationMockTokenParser) DeleteUser(ctx context.Context, firebaseUserID string) error {
	return nil
}
