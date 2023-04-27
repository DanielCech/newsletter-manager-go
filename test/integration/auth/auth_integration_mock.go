package auth

import (
	"context"
	"event-facematch-backend/firebase"
	_ "event-facematch-backend/types"
	id "event-facematch-backend/types/id"
)

type IntegrationMockTokenParser struct {
	AuthorIDMappings map[string]id.Author
}

func NewIntegrationMockTokenParser() IntegrationMockTokenParser {
	return IntegrationMockTokenParser{AuthorIDMappings: make(map[string]id.Author)}
}

func (p *IntegrationMockTokenParser) VerifyIDToken(ctx context.Context, idToken string) (firebase.VerifiedToken, error) {
	AuthorID, ok := p.AuthorIDMappings[idToken]
	if !ok {
		return firebase.VerifiedToken{
			AuthorID:         nil,
			FirebaseAuthorID: idToken,
		}, nil
	}

	return firebase.VerifiedToken{
		AuthorID:         &AuthorID,
		FirebaseAuthorID: idToken,
	}, nil
}

func (p *IntegrationMockTokenParser) SetAuthorIDCustomClaim(ctx context.Context, firebaseAuthorID string, AuthorID id.Author) error {
	p.AuthorIDMappings[firebaseAuthorID] = AuthorID
	return nil
}

func (p *IntegrationMockTokenParser) DeleteUser(ctx context.Context, firebaseAuthorID string) error {
	return nil
}
