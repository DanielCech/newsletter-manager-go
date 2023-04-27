package id

import (
	"github.com/google/uuid"
)

//go:generate go run go.strv.io/tea/cmd/tea gen id -i ./id.go -o ./id_gen.go

type (
	Author       uuid.UUID
	Newsletter   uuid.UUID
	Email        uuid.UUID
	Subscription uuid.UUID
)

type (
	RefreshToken string
)
