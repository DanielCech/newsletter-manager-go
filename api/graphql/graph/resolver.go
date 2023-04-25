package graph

import (
	"errors"
)

//go:generate go run github.com/99designs/gqlgen generate

var errMissingUserID = errors.New("missing user id")

type Resolver struct {
	userService    UserService
	sessionService SessionService
}

func NewResolver(userService UserService, sessionService SessionService) *Resolver {
	return &Resolver{
		userService:    userService,
		sessionService: sessionService,
	}
}
