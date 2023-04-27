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

func NewAuthor() Author {
	return Author(uuid.New())
}

func (i Author) String() string {
	return uuid.UUID(i).String()
}

func (i Author) Empty() bool {
	return uuid.UUID(i) == uuid.Nil
}

func (i Author) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *Author) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "Author", data)
}

func (i *Author) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "Author", data)
}

func NewNewsletter() Newsletter {
	return Newsletter(uuid.New())
}

func (i Newsletter) String() string {
	return uuid.UUID(i).String()
}

func (i Newsletter) Empty() bool {
	return uuid.UUID(i) == uuid.Nil
}

func (i Newsletter) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *Newsletter) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "Newsletter", data)
}

func (i *Newsletter) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "Newsletter", data)
}

func NewEmail() Email {
	return Email(uuid.New())
}

func (i Email) String() string {
	return uuid.UUID(i).String()
}

func (i Email) Empty() bool {
	return uuid.UUID(i) == uuid.Nil
}

func (i Email) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *Email) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "Email", data)
}

func (i *Email) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "Email", data)
}

func NewSubscription() Subscription {
	return Subscription(uuid.New())
}

func (i Subscription) String() string {
	return uuid.UUID(i).String()
}

func (i Subscription) Empty() bool {
	return uuid.UUID(i) == uuid.Nil
}

func (i Subscription) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(i).String()), nil
}

func (i *Subscription) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(i), "Subscription", data)
}

func (i *Subscription) Scan(data any) error {
	return scanUUID((*uuid.UUID)(i), "Subscription", data)
}

func (i *RefreshToken) UnmarshalText(data []byte) error {
	*i = RefreshToken(data)
	return nil
}

func (i RefreshToken) MarshalText() ([]byte, error) {
	return []byte(i), nil
}
