package id

import (
	"github.com/google/uuid"
)

//go:generate go run go.strv.io/tea/cmd/tea gen id -i ./id.go -o ./id_gen.go

type (
	User uuid.UUID
)

type (
	RefreshToken string
)
