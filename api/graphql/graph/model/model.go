package model

import (
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// MarshalUUID marshals uuid.UUID to string.
func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(u.String())
}

// UnmarshalUUID unmarshals value into the uuid.UUID.
func UnmarshalUUID(v any) (uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("invalid UUID: %w", err)
		}
		return u, nil
	default:
		return uuid.UUID{}, fmt.Errorf("%T must be a string", v)
	}
}

// MarshalDateTime marshals time.Time to string.
func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalString(t.UTC().Format(time.RFC3339))
}

// UnmarshalDateTime unmarshals value into the time.Time.
func UnmarshalDateTime(v any) (time.Time, error) {
	switch v := v.(type) {
	case string:
		parsedTime, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid DateTime: %w", err)
		}
		return parsedTime, nil
	default:
		return time.Time{}, fmt.Errorf("%T must be a string", v)
	}
}
