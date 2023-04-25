package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type SecretConfiger interface {
	LocalPath() *string
	SecretARN() *string
}

type Secret[T any] struct {
	LocalPath *string
	SecretARN *string
}

func NewSecret[T any](s SecretConfiger) Secret[T] {
	return Secret[T]{
		LocalPath: s.LocalPath(),
		SecretARN: s.SecretARN(),
	}
}

func (s Secret[T]) Resolve(ctx context.Context, sm SecretsManager) (*T, error) {
	t := new(T)

	switch {
	case s.LocalPath != nil:
		b, err := os.ReadFile(*s.LocalPath)
		if err != nil {
			return nil, fmt.Errorf("reading local secret value: %w", err)
		}
		if err = json.Unmarshal(b, t); err != nil {
			return nil, fmt.Errorf("parsing local secret value: %w", err)
		}
	case s.SecretARN != nil:
		if err := sm.GetAndParseSecretValue(ctx, *s.SecretARN, t); err != nil {
			return nil, fmt.Errorf("creating config: %w", err)
		}
	default:
		return nil, ErrMissingPathOrARN
	}

	if err := validate.Struct(t); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	return t, nil
}
