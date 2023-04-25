package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"strv-template-backend-go-api/types/id"
	"strv-template-backend-go-api/util/timesource"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	refreshTokenLen = 16
)

// Factory contains dependencies that are needed for sessions creation.
type Factory struct {
	secret                 []byte
	timeSource             timesource.TimeSource
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

// NewFactory returns new instance of session Factory.
func NewFactory(
	secret []byte,
	timeSource timesource.TimeSource,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
) (Factory, error) {
	if err := newFactoryValidate(secret, timeSource, accessTokenExpiration, refreshTokenExpiration); err != nil {
		return Factory{}, err
	}
	return Factory{
		secret:                 secret,
		timeSource:             timeSource,
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

// NewAccessToken returns new instance of AccessToken.
func (f Factory) NewAccessToken(customClaims Claims) (*AccessToken, error) {
	now := f.timeSource.Now()
	expiration := now.Add(f.accessTokenExpiration)
	signedData, err := f.signAccessToken(customClaims, now, expiration)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}
	return &AccessToken{
		timeSource: f.timeSource,
		SignedData: signedData,
		ExpiresAt:  expiration,
	}, nil
}

// ParseAccessToken returns new instance of AccessToken based on provided signed data.
func (f Factory) ParseAccessToken(data string) (*AccessToken, error) {
	tokenClaims, err := f.parseAccessToken(data)
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.Parse(tokenClaims.Subject)
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		timeSource: f.timeSource,
		Claims: Claims{
			UserID: id.User(userUUID),
			Custom: tokenClaims.Custom,
		},
		SignedData: data,
		ExpiresAt:  tokenClaims.ExpiresAt.Time,
	}, nil
}

// NewRefreshToken returns new instance of RefreshToken.
func (f Factory) NewRefreshToken(userID id.User) (*RefreshToken, error) {
	data := make([]byte, refreshTokenLen)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	now := f.timeSource.Now()
	expiresAt := now.Add(f.refreshTokenExpiration)
	return &RefreshToken{
		timeSource: f.timeSource,
		ID:         id.RefreshToken(base64.StdEncoding.EncodeToString(data)),
		UserID:     userID,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
	}, nil
}

// NewRefreshTokenFromFields returns new instance of RefreshToken based on existing fields.
// This can be useful for repositories when converting results from repository to domain models based on consistent data.
// There is no validity check, it is responsibility of caller to ensure all fields are correct.
func (f Factory) NewRefreshTokenFromFields(
	id id.RefreshToken,
	userID id.User,
	expiresAt time.Time,
	createdAt time.Time,
) *RefreshToken {
	refreshToken := &RefreshToken{
		timeSource: f.timeSource,
		ID:         id,
		UserID:     userID,
		ExpiresAt:  expiresAt,
		CreatedAt:  createdAt,
	}
	return refreshToken
}

// NewSession returns new instance of session.
// New access and refresh token are created based on provided custom jwtClaims.
func (f Factory) NewSession(claims Claims) (*Session, error) {
	accessToken, err := f.NewAccessToken(claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := f.NewRefreshToken(claims.UserID)
	if err != nil {
		return nil, err
	}
	return &Session{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

// parseAccessToken parses signed data and returns jwtClaims according to RFC 7519.
// Parsable are only those tokens which are signed by HMAC method.
func (f Factory) parseAccessToken(data string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(data, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt signing method")
		}
		return f.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing jwt token with claims: %w", err)
	}

	c, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return c, nil
}

// signAccessToken returns signed data of access token.
// As a signing method is used HMAC with SHA256.
func (f Factory) signAccessToken(claims Claims, issuedAt, expiresAt time.Time) (string, error) {
	c := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Custom: claims.Custom,
	}

	signedData, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &c).SignedString(f.secret)
	if err != nil {
		return "", fmt.Errorf("new jwt with claims: %w", err)
	}

	return signedData, nil
}

func newFactoryValidate(
	secret []byte,
	timeSource timesource.TimeSource,
	accessTokenExpiration time.Duration,
	refreshTokenExpiration time.Duration,
) error {
	if len(secret) == 0 {
		return ErrInvalidSecret
	}
	if timeSource == nil {
		return ErrInvalidTimeSource
	}
	if accessTokenExpiration == 0 {
		return errors.New("invalid access token expiration")
	}
	if refreshTokenExpiration == 0 {
		return errors.New("invalid refresh token expiration")
	}
	return nil
}

// jwtClaims contains registered jwtClaims with custom ones.
type jwtClaims struct {
	jwt.RegisteredClaims
	Custom CustomClaims `json:"custom"`
}
