package common

import (
	"context"
	"newsletter-manager-go/test/integration/generate/swagger"
)

type TestingAuthor struct {
	Context  context.Context
	AuthorID string
	Session  *swagger.Session
}

func NewUser() TestingAuthor {
	return TestingAuthor{
		Context: context.Background(),
	}
}

func (author *TestingAuthor) UpdateWithResponse(authorID string, session *swagger.Session) {
	author.AuthorID = authorID
	author.Session = session
	author.Context = context.WithValue(author.Context, swagger.ContextAccessToken, session.AccessToken)
}

func (author *TestingAuthor) UpdateWithSession(session *swagger.Session) {
	author.Session = session
}
