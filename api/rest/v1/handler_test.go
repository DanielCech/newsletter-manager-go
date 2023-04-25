package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	errTest = errors.New("test error")
)

type mockTokenParser struct {
	mock.Mock
}

func (m *mockTokenParser) ParseAccessToken(token string) (*domsession.AccessToken, error) {
	args := m.Called(token)
	return args.Get(0).(*domsession.AccessToken), args.Error(1)
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Create(_ context.Context, createUserInput domuser.CreateUserInput) (*domuser.User, *domsession.Session, error) {
	args := m.Called(createUserInput)
	return args.Get(0).(*domuser.User), args.Get(1).(*domsession.Session), args.Error(2)
}

func (m *mockUserService) Read(_ context.Context, userID id.User) (*domuser.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockUserService) ReadByCredentials(_ context.Context, email types.Email, password types.Password) (*domuser.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockUserService) ChangePassword(_ context.Context, userID id.User, oldPassword, newPassword types.Password) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *mockUserService) List(_ context.Context) ([]domuser.User, error) {
	args := m.Called()
	return args.Get(0).([]domuser.User), args.Error(1)
}

type mockSessionService struct {
	mock.Mock
}

func (m *mockSessionService) Create(_ context.Context, email types.Email, password types.Password) (*domsession.Session, *domuser.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*domsession.Session), args.Get(1).(*domuser.User), args.Error(2)
}

func (m *mockSessionService) Destroy(_ context.Context, refreshTokenID id.RefreshToken) error {
	args := m.Called(refreshTokenID)
	return args.Error(0)
}

func (m *mockSessionService) Refresh(_ context.Context, refreshTokenID id.RefreshToken) (*domsession.Session, error) {
	args := m.Called(refreshTokenID)
	return args.Get(0).(*domsession.Session), args.Error(1)
}

type mocks struct {
	tokenParser    *mockTokenParser
	userService    *mockUserService
	sessionService *mockSessionService
}

func newMocks() mocks {
	return mocks{
		tokenParser:    &mockTokenParser{},
		userService:    &mockUserService{},
		sessionService: &mockSessionService{},
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.tokenParser.AssertExpectations(t)
	m.userService.AssertExpectations(t)
	m.sessionService.AssertExpectations(t)
}

func assertResponseBody(t *testing.T, expectedBody any, body *bytes.Buffer) {
	t.Helper()

	if expectedBody == http.NoBody {
		assert.Empty(t, body)
		return
	}

	bodyData := body.Bytes()
	expectedBodyData, err := json.Marshal(expectedBody)
	assert.NoError(t, err)
	assert.Equal(t, bytes.TrimSpace(expectedBodyData), bytes.TrimSpace(bodyData))
}
