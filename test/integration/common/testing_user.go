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

func NewUser(email string, password string) TestingAuthor {
	return TestingAuthor{
		Context: context.Background(),
	}
}

func (author *TestingAuthor) UpdateWith(resp swagger.CreateAuthorResp) {
	author.AuthorID = resp.Author.Id
	author.Session = resp.Session
}
