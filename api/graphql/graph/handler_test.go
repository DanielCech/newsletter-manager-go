package graph

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"strv-template-backend-go-api/api/graphql/middleware"
	"strv-template-backend-go-api/database/sql"
	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) HashPassword(password []byte) ([]byte, error) {
	args := m.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockHasher) CompareHashAndPassword(hash, password []byte) bool {
	args := m.Called(hash, password)
	return args.Bool(0)
}

type mockTimeSource struct {
	mock.Mock
}

func (m *mockTimeSource) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Create(_ context.Context, createUserInput domuser.CreateUserInput) (*domuser.User, *domsession.Session, error) {
	args := m.Called(createUserInput)
	return args.Get(0).(*domuser.User), args.Get(1).(*domsession.Session), args.Error(2)
}

func (m *mockUserService) Read(_ context.Context, id id.User) (*domuser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockUserService) ReadByCredentials(_ context.Context, email types.Email, password types.Password) (*domuser.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*domuser.User), args.Error(1)
}

func (m *mockUserService) ChangePassword(_ context.Context, id id.User, oldPassword, newPassword types.Password) error {
	args := m.Called(id, oldPassword, newPassword)
	return args.Error(0)
}

func (m *mockUserService) List(context.Context) ([]domuser.User, error) {
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

type mockTokenParser struct{}

func (m *mockTokenParser) ParseAccessToken(string) (*domsession.AccessToken, error) {
	return nil, nil
}

type dbContext struct {
	context.Context
}

func (d dbContext) Ctx() context.Context {
	return d.Context
}

type mockDataSource struct {
	mock.Mock
}

func (m *mockDataSource) AcquireConnCtx(context.Context) (sql.DataContext, error) {
	args := m.Called()
	return args.Get(0).(sql.DataContext), args.Error(1)
}

func (m *mockDataSource) ReleaseConnCtx(dctx sql.DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

func (m *mockDataSource) Begin(context.Context) (sql.DataContext, error) {
	args := m.Called()
	return args.Get(0).(sql.DataContext), args.Error(1)
}

func (m *mockDataSource) Commit(dctx sql.DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

func (m *mockDataSource) Rollback(dctx sql.DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

type mocks struct {
	hasher         *mockHasher
	timeSource     *mockTimeSource
	userService    *mockUserService
	sessionService *mockSessionService
	dataSource     *mockDataSource
}

func newMocks() mocks {
	return mocks{
		hasher:         &mockHasher{},
		timeSource:     &mockTimeSource{},
		userService:    &mockUserService{},
		sessionService: &mockSessionService{},
		dataSource:     &mockDataSource{},
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.hasher.AssertExpectations(t)
	m.timeSource.AssertExpectations(t)
	m.userService.AssertExpectations(t)
	m.sessionService.AssertExpectations(t)
	m.dataSource.AssertExpectations(t)
}

var queryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	if strings.Compare(expectedSQL, actualSQL) == 0 {
		return nil
	}
	return fmt.Errorf(`could not match actual sql: "%s" with expected: "%s"`, actualSQL, expectedSQL)
})

func Test_NewController(t *testing.T) {
	controller, err := NewController(
		&mockUserService{},
		&mockSessionService{},
		&mockTokenParser{},
		&mockDataSource{},
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
		&mockSessionService{},
		&mockTokenParser{},
		&mockDataSource{},
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
		&mockUserService{},
		&mockSessionService{},
		&mockTokenParser{},
		&mockDataSource{},
		middleware.CORSConfig{
			AllowedCredentials: ptr(true),
			MaxAge:             ptr(300),
		},
		zap.NewNop(),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, controller.Mux)
}

func Test_newControllerValidate(t *testing.T) {
	err := newControllerValidate(&mockUserService{}, &mockSessionService{}, &mockTokenParser{}, &mockDataSource{}, zap.NewNop())
	assert.NoError(t, err)

	err = newControllerValidate(&mockUserService{}, &mockSessionService{}, &mockTokenParser{}, &mockDataSource{}, nil)
	assert.EqualError(t, err, "invalid logger")

	err = newControllerValidate(&mockUserService{}, &mockSessionService{}, &mockTokenParser{}, nil, nil)
	assert.EqualError(t, err, "invalid data source")

	err = newControllerValidate(&mockUserService{}, &mockSessionService{}, nil, nil, nil)
	assert.EqualError(t, err, "invalid token parser")

	err = newControllerValidate(&mockUserService{}, nil, nil, nil, nil)
	assert.EqualError(t, err, "invalid session service")

	err = newControllerValidate(nil, nil, nil, nil, nil)
	assert.EqualError(t, err, "invalid user service")
}
