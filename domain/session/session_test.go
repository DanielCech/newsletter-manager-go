package session

import (
	"testing"
	"time"

	"newsletter-manager-go/domain/user"
	"newsletter-manager-go/types/id"

	"github.com/stretchr/testify/assert"
)

func Test_NewClaims(t *testing.T) {
	authorID := id.NewUser()
	role := user.RoleUser
	expected := Claims{
		AuthorID: authorID,
		Custom:   CustomClaims{UserRole: user.RoleUser},
	}
	customClaims, err := NewClaims(authorID, role)
	assert.NoError(t, err)
	assert.Equal(t, expected, customClaims)

	customClaims, err = NewClaims(id.Author{}, role)
	assert.Error(t, err)
	assert.Empty(t, customClaims)
}

func Test_Claims_Valid(t *testing.T) {
	customClaims := Claims{
		AuthorID: id.NewUser(),
		Custom:   CustomClaims{UserRole: user.RoleUser},
	}
	err := customClaims.Valid()
	assert.NoError(t, err)

	customClaims.Custom.UserRole = ""
	err = customClaims.Valid()
	assert.Error(t, err)

	customClaims.AuthorID = id.Author{}
	err = customClaims.Valid()
	assert.ErrorIs(t, err, user.ErrInvalidAuthorID)
}

func Test_AccessToken_IsExpired(t *testing.T) {
	timeSource := &mockTimeSource{}
	now := time.Now()
	timeSource.On("Now").Return(now)
	accessToken := AccessToken{
		timeSource: timeSource,
		ExpiresAt:  now.Add(time.Minute),
	}
	assert.False(t, accessToken.IsExpired())

	accessToken.ExpiresAt = now.Add(-time.Hour)
	assert.True(t, accessToken.IsExpired())

	timeSource.AssertExpectations(t)
}

func Test_RefreshToken_IsExpired(t *testing.T) {
	timeSource := &mockTimeSource{}
	now := time.Now()
	timeSource.On("Now").Return(now)
	refreshToken := RefreshToken{
		timeSource: timeSource,
		ExpiresAt:  now.Add(time.Minute),
	}
	assert.False(t, refreshToken.IsExpired())

	refreshToken.ExpiresAt = now.Add(-time.Hour)
	assert.True(t, refreshToken.IsExpired())

	timeSource.AssertExpectations(t)
}
