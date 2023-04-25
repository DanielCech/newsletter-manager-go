package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"strv-template-backend-go-api/api/rest/middleware"
	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

type mockUserService struct{}

func (mockUserService) Create(context.Context, domuser.CreateUserInput) (*domuser.User, *domsession.Session, error) {
	return nil, nil, nil
}

func (mockUserService) Read(context.Context, id.User) (*domuser.User, error) {
	return nil, nil
}

func (mockUserService) ReadByCredentials(context.Context, types.Email, types.Password) (*domuser.User, error) {
	return nil, nil
}

func (mockUserService) ChangePassword(context.Context, id.User, types.Password, types.Password) error {
	return nil
}

func (mockUserService) List(context.Context) ([]domuser.User, error) {
	return nil, nil
}

type mockSessionService struct{}

func (mockSessionService) Create(context.Context, types.Email, types.Password) (*domsession.Session, *domuser.User, error) {
	return nil, nil, nil
}

func (mockSessionService) Destroy(context.Context, id.RefreshToken) error {
	return nil
}

func (mockSessionService) Refresh(context.Context, id.RefreshToken) (*domsession.Session, error) {
	return nil, nil
}

type mockTokenParser struct{}

func (mockTokenParser) ParseAccessToken(string) (*domsession.AccessToken, error) {
	return nil, nil
}

func Test_NewController(t *testing.T) {
	controller, err := NewController(
		mockUserService{},
		mockSessionService{},
		mockTokenParser{},
		middleware.CORSConfig{
			AllowedCredentials: ptr(true),
			MaxAge:             ptr(300),
		},
		zap.NewNop(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, controller.userService)
	assert.NotNil(t, controller.sessionService)
	assert.NotNil(t, controller.tokenParser)
	assert.NotNil(t, controller.logger)

	controller, err = NewController(
		nil,
		mockSessionService{},
		mockTokenParser{},
		middleware.CORSConfig{
			AllowedCredentials: ptr(true),
			MaxAge:             ptr(300),
		},
		zap.NewNop(),
	)
	assert.Error(t, err)
	assert.Empty(t, controller)
}

func Test_Controller_initRouter(t *testing.T) {
	controller, err := NewController(
		mockUserService{},
		mockSessionService{},
		mockTokenParser{},
		middleware.CORSConfig{
			AllowedCredentials: ptr(true),
			MaxAge:             ptr(300),
		},
		zap.NewNop(),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, controller.Mux)
}

func Test_Controller_OpenAPI(t *testing.T) {
	controller, err := NewController(
		mockUserService{},
		mockSessionService{},
		mockTokenParser{},
		middleware.CORSConfig{
			AllowedCredentials: ptr(true),
			MaxAge:             ptr(300),
		},
		zap.NewNop(),
	)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/api/openapi.yaml", http.NoBody)
	require.NoError(t, err)

	controller.OpenAPI(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func Test_newControllerValidate(t *testing.T) {
	err := newControllerValidate(mockUserService{}, mockSessionService{}, mockTokenParser{}, zap.NewNop())
	assert.NoError(t, err)

	err = newControllerValidate(mockUserService{}, mockSessionService{}, mockTokenParser{}, nil)
	assert.EqualError(t, err, "invalid logger")

	err = newControllerValidate(mockUserService{}, mockSessionService{}, nil, nil)
	assert.EqualError(t, err, "invalid token parser")

	err = newControllerValidate(mockUserService{}, nil, nil, nil)
	assert.EqualError(t, err, "invalid session service")

	err = newControllerValidate(nil, nil, nil, nil)
	assert.EqualError(t, err, "invalid user service")
}
