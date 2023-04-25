package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/cors"
)

// ArrayWithTextUnmarshaller is able to unmarshal array in string representation.
// Example: VALUE1,VALUE2,VALUE3
type ArrayWithTextUnmarshaller []string

func (a *ArrayWithTextUnmarshaller) UnmarshalText(text []byte) error {
	if a == nil {
		return errors.New("unmarshal text: nil pointer")
	}
	*a = strings.Split(string(text), ",")
	return nil
}

// CORSConfig contains required fields for CORS configuration.
type CORSConfig struct {
	AllowedOrigins     ArrayWithTextUnmarshaller `json:"allowed_origins" yaml:"allowed_origins" env:"CORS_ALLOWED_ORIGINS" validate:"required"`
	AllowedMethods     ArrayWithTextUnmarshaller `json:"allowed_methods" yaml:"allowed_methods" env:"CORS_ALLOWED_METHODS" validate:"required"`
	AllowedHeaders     ArrayWithTextUnmarshaller `json:"allowed_headers" yaml:"allowed_headers" env:"CORS_ALLOWED_HEADERS" validate:"required"`
	AllowedCredentials *bool                     `json:"allowed_credentials" yaml:"allowed_credentials" env:"CORS_ALLOWED_CREDENTIALS" validate:"required"`
	MaxAge             *int                      `json:"max_age" yaml:"max_age" env:"CORS_MAX_AGE" validate:"required"`
}

func NewCORSHandler(cfg CORSConfig) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: *cfg.AllowedCredentials,
		MaxAge:           *cfg.MaxAge,
	})
}
