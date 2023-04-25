package id

import (
	"fmt"

	"github.com/google/uuid"
)

func unmarshalUUID(u *uuid.UUID, idTypeName string, data []byte) error {
	if err := u.UnmarshalText(data); err != nil {
		return fmt.Errorf("parsing %q id value: %w", idTypeName, err)
	}
	return nil
}

func scanUUID(u *uuid.UUID, idTypeName string, data any) error {
	if err := u.Scan(data); err != nil {
		return fmt.Errorf("scanning %q id value: %w", idTypeName, err)
	}
	return nil
}

func NewUser() User {
	return User(uuid.New())
}

func (i User) String() string {
	return uuid.UUID(i).String()
}

func (i User) Empty() bool {
	return uuid.UUID(i) == uuid.Nil
}

func (i User) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *User) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "User", data)
}

func (i *User) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "User", data)
}

func (i *RefreshToken) UnmarshalText(data []byte) error {
	*i = RefreshToken(data)
	return nil
}

func (i RefreshToken) MarshalText() ([]byte, error) {
	return []byte(i), nil
}
