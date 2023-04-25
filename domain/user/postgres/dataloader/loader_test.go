package dataloader

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"strv-template-backend-go-api/database/sql"
	"strv-template-backend-go-api/domain/user/postgres/dataloader/query"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

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
	dataSource *mockDataSource
	querier    pgxmock.PgxPoolIface
}

func newMocks(t *testing.T) mocks {
	t.Helper()
	querier, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(queryMatcher))
	require.NoError(t, err)
	return mocks{
		dataSource: &mockDataSource{},
		querier:    querier,
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.dataSource.AssertExpectations(t)
	assert.NoError(t, m.querier.ExpectationsWereMet())
}

var queryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	if strings.Compare(expectedSQL, actualSQL) == 0 {
		return nil
	}
	return fmt.Errorf(`could not match actual sql: "%s" with expected: "%s"`, actualSQL, expectedSQL)
})

func Test_New(t *testing.T) {
	loader := New(&mockDataSource{})
	assert.NotEmpty(t, loader.User)
}

func Test_WithLoader(t *testing.T) {
	expected := New(&mockDataSource{})
	ctx := WithLoader(context.Background(), expected)
	loader, ok := ctx.Value(contextKey.loader).(*Loader)
	assert.True(t, ok)
	assert.Equal(t, expected, loader)
}

func Test_LoaderFromCtx(t *testing.T) {
	expected := New(&mockDataSource{})
	ctx := context.WithValue(context.Background(), contextKey.loader, expected)
	loader, ok := LoaderFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, loader)
}

func Test_ReadUser(t *testing.T) {
	ctx := context.Background()
	user := newUser()
	mocks := newMocks(t)
	dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
	mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
	mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
	queryArgs := pgx.NamedArgs{
		"ids": []string{user.ID.String()},
	}
	rows := pgxmock.NewRows([]string{"id", "referrer_id", "name", "email", "password_hash", "role", "created_at", "updated_at"})
	rows = rows.AddRow(user.ID, user.ReferrerID, user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	mocks.querier.ExpectQuery(query.ListUsersByIDs).WithArgs(queryArgs).WillReturnRows(rows)
	loader := New(mocks.dataSource)
	ctx = WithLoader(ctx, loader)

	result, err := ReadUser(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mocks.assertExpectations(t)

	result, err = ReadUser(context.Background(), user.ID)
	assert.EqualError(t, err, "missing data loader")
	assert.Empty(t, result)
}
