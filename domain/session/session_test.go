package session

import (
	"testing"
	"time"

	"strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"

	"github.com/stretchr/testify/assert"
)

func Test_NewClaims(t *testing.T) {
	userID := id.NewUser()
	role := user.RoleUser
	expected := Claims{
		UserID: userID,
		Custom: CustomClaims{UserRole: user.RoleUser},
	}
	customClaims, err := NewClaims(userID, role)
	assert.NoError(t, err)
	assert.Equal(t, expected, customClaims)

	customClaims, err = NewClaims(id.User{}, role)
	assert.Error(t, err)
	assert.Empty(t, customClaims)
}

func Test_Claims_Valid(t *testing.T) {
	customClaims := Claims{
		UserID: id.NewUser(),
		Custom: CustomClaims{UserRole: user.RoleUser},
	}
	err := customClaims.Valid()
	assert.NoError(t, err)

	customClaims.Custom.UserRole = ""
	err = customClaims.Valid()
	assert.Error(t, err)

	customClaims.UserID = id.User{}
	err = customClaims.Valid()
	assert.ErrorIs(t, err, user.ErrInvalidUserID)
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
