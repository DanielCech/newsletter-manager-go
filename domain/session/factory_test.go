package session

import (
	"testing"
	"time"

	"strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testSecret = []byte("45asd4a5d4a3ds48a")
)

type mockTimeSource struct {
	mock.Mock
}

func (m *mockTimeSource) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func Test_NewFactory(t *testing.T) {
	timeSource := &mockTimeSource{}
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour
	expected := Factory{
		secret:                 testSecret,
		timeSource:             timeSource,
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
	factory, err := NewFactory(testSecret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	assert.NoError(t, err)
	assert.Equal(t, expected, factory)

	factory, err = NewFactory(nil, timeSource, accessTokenExpiration, refreshTokenExpiration)
	assert.Error(t, err)
	assert.Empty(t, factory)
}

func Test_Factory_NewAccessToken(t *testing.T) {
	now := time.Now()
	timeSource := &mockTimeSource{}
	timeSource.On("Now").Return(now)
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour
	factory, err := NewFactory(testSecret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	require.NoError(t, err)
	claims := Claims{
		UserID: id.NewUser(),
		Custom: CustomClaims{UserRole: user.RoleUser},
	}

	accessToken, err := factory.NewAccessToken(claims)
	assert.NoError(t, err)
	assert.Equal(t, timeSource, accessToken.timeSource)
	assert.NotEmpty(t, accessToken.SignedData)
	assert.True(t, accessToken.ExpiresAt.Equal(now.Add(accessTokenExpiration)))

	timeSource.AssertExpectations(t)
}

func Test_Factory_ParseAccessToken(t *testing.T) {
	data := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwNzA5OTg4MDAsImlhdCI6MTY3NTc4Mjg1MCwianRpIjoiNGU2NmNiMjQtYTM3Ni00YzE2LWFkYmMtZjliZTE0ZjU0YWJjIiwic3ViIjoiZjY1MGUyNmMtYmNlMC00OGIyLTllYTQtMTI2MTkzZDM0MWUwIiwiY3VzdG9tIjp7InVzZXJfcm9sZSI6InVzZXIifX0.cNKMRHuH2hGeaaxhi6yf8sJObSk7AgNNkW-E_44lXuM"
	testTime := time.Unix(4070995200, 0)
	timeSource := &mockTimeSource{}
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour
	factory, err := NewFactory(testSecret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	require.NoError(t, err)
	expectedClaims := Claims{
		UserID: id.User(uuid.MustParse("f650e26c-bce0-48b2-9ea4-126193d341e0")),
		Custom: CustomClaims{UserRole: user.RoleUser},
	}

	accessToken, err := factory.ParseAccessToken(data)
	assert.NoError(t, err)
	assert.Equal(t, timeSource, accessToken.timeSource)
	assert.Equal(t, expectedClaims, accessToken.Claims)
	assert.Equal(t, data, accessToken.SignedData)
	assert.True(t, testTime.Add(accessTokenExpiration).Equal(accessToken.ExpiresAt))

	data = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI2MjEyMDAsIkN1c3RvbSI6eyJVc2VySUQiOiI0MWI0NDUxMy0zYTFjLTRlNjEtOTFhMy0wNTc5ZDk1YTA2MmEiLCJVc2VyUm9sZSI6InVzZXIifX0.ZA9TdAks1EJg-LDKu44IBORT0EYxtQ4xUMDAp7M5oCQ"
	accessToken, err = factory.ParseAccessToken(data)
	assert.Error(t, err)
	assert.Empty(t, accessToken)

	data = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwNzA5OTUyMDAsImlhdCI6MTY3MjYxNzYwMCwianRpIjoiZjY1MGUyNmMtYmNlMC00OGIyLTllYTQtMTI2MTkzZDM0MWUwIiwic3ViIjoiZjY1MGUyIiwiY3VzdG9tIjp7InVzZXJfcm9sZSI6InVzZXIifX0.gWj6NjEP1pxM4ELIV14XJ_EKlCNu4UrTCIoeInxMbLM"
	accessToken, err = factory.ParseAccessToken(data)
	assert.Error(t, err)
	assert.Empty(t, accessToken)
}

func Test_Factory_NewRefreshToken(t *testing.T) {
	now := time.Now()
	timeSource := &mockTimeSource{}
	timeSource.On("Now").Return(now)
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour
	factory, err := NewFactory(testSecret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	require.NoError(t, err)
	userID := id.NewUser()

	refreshToken, err := factory.NewRefreshToken(userID)
	assert.NoError(t, err)
	assert.Equal(t, timeSource, refreshToken.timeSource)
	assert.NotEmpty(t, refreshToken.ID)
	assert.Equal(t, userID, refreshToken.UserID)
	assert.True(t, refreshToken.ExpiresAt.Equal(now.Add(refreshTokenExpiration)))
	assert.True(t, refreshToken.CreatedAt.Equal(now))
	timeSource.AssertExpectations(t)
}

func Test_Factory_NewRefreshTokenFromFields(t *testing.T) {
	now := time.Now()
	timeSource := &mockTimeSource{}
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour
	factory, err := NewFactory(testSecret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	require.NoError(t, err)
	refreshTokenID := id.RefreshToken("5asd4a6d4a36d45as36da")
	userID := id.User(uuid.MustParse("4049b2e8-767e-475e-b78d-6baaf356cdc1"))
	expiresAt := now.Add(refreshTokenExpiration)
	createdAt := now

	refreshToken := factory.NewRefreshTokenFromFields(refreshTokenID, userID, expiresAt, createdAt)
	assert.Equal(t, timeSource, refreshToken.timeSource)
	assert.Equal(t, refreshTokenID, refreshToken.ID)
	assert.True(t, refreshToken.ExpiresAt.Equal(expiresAt))
	assert.True(t, refreshToken.CreatedAt.Equal(createdAt))
}

func Test_Factory_NewSession(t *testing.T) {
	timeSource := &mockTimeSource{}
	timeSource.On("Now").Return(time.Now())
	factory, err := NewFactory(testSecret, timeSource, time.Hour, time.Hour)
	require.NoError(t, err)
	claims := Claims{
		UserID: id.NewUser(),
		Custom: CustomClaims{UserRole: user.RoleUser},
	}

	session, err := factory.NewSession(claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, session.AccessToken)
	assert.NotEmpty(t, session.RefreshToken)

	timeSource.AssertExpectations(t)
}

func Test_parseAccessToken(t *testing.T) {
	data := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwNzA5OTUyMDAsImlhdCI6MTY3MjYxNzYwMCwianRpIjoiZjY1MGUyNmMtYmNlMC00OGIyLTllYTQtMTI2MTkzZDM0MWUwIiwic3ViIjoiZjY1MGUyNmMtYmNlMC00OGIyLTllYTQtMTI2MTkzZDM0MWUxIiwiY3VzdG9tIjp7InVzZXJfcm9sZSI6InVzZXIifX0.wTxt8t_F-KDZlWH8xiGns1z2g4YHGQAXxSedegWIGOU"
	factory, err := NewFactory(testSecret, &mockTimeSource{}, time.Hour, time.Hour)
	require.NoError(t, err)
	expiresAt := time.Unix(4070995200, 0)
	issuedAt := time.Unix(1672617600, 0)
	require.NoError(t, err)
	expected := &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "f650e26c-bce0-48b2-9ea4-126193d341e1",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        "f650e26c-bce0-48b2-9ea4-126193d341e0",
		},
		Custom: CustomClaims{
			UserRole: user.RoleUser,
		},
	}

	claims, err := factory.parseAccessToken(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, claims)
}

func Test_signAccessToken(t *testing.T) {
	factory, err := NewFactory(testSecret, &mockTimeSource{}, time.Hour, time.Hour)
	require.NoError(t, err)
	claims := Claims{
		UserID: id.User(uuid.MustParse("f650e26c-bce0-48b2-9ea4-126193d341e1")),
		Custom: CustomClaims{UserRole: user.RoleUser},
	}
	issuedAt := time.Unix(4070994200, 0)
	expiresAt := time.Unix(4070995200, 0)

	signedData, err := factory.signAccessToken(claims, issuedAt, expiresAt)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmNjUwZTI2Yy1iY2UwLTQ4YjItOWVhNC0xMjYxOTNkMzQxZTEiLCJleHAiOjQwNzA5OTUyMDAsImlhdCI6NDA3MDk5NDIwMCwiY3VzdG9tIjp7InVzZXJfcm9sZSI6InVzZXIifX0.b_RCNZVpPqeKxG3yv-8q3wqO-7jssoZConTeR69YayA", signedData)
}

func Test_newFactoryValidate(t *testing.T) {
	secret := testSecret
	timeSource := &mockTimeSource{}
	accessTokenExpiration := time.Hour
	refreshTokenExpiration := time.Hour

	err := newFactoryValidate(secret, timeSource, accessTokenExpiration, refreshTokenExpiration)
	assert.NoError(t, err)

	err = newFactoryValidate(secret, timeSource, accessTokenExpiration, time.Duration(0))
	assert.EqualError(t, err, "invalid refresh token expiration")

	err = newFactoryValidate(secret, timeSource, time.Duration(0), time.Duration(0))
	assert.EqualError(t, err, "invalid access token expiration")

	err = newFactoryValidate(secret, nil, time.Duration(0), time.Duration(0))
	assert.ErrorIs(t, err, ErrInvalidTimeSource)

	err = newFactoryValidate(nil, nil, time.Duration(0), time.Duration(0))
	assert.ErrorIs(t, err, ErrInvalidSecret)
}
