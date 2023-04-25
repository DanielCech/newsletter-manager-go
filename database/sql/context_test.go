package sql

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDataSource struct {
	mock.Mock
}

func (m *mockDataSource) AcquireConnCtx(context.Context) (DataContext, error) {
	args := m.Called()
	return args.Get(0).(DataContext), args.Error(1)
}

func (m *mockDataSource) ReleaseConnCtx(dctx DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

func (m *mockDataSource) Begin(context.Context) (DataContext, error) {
	args := m.Called()
	return args.Get(0).(DataContext), args.Error(1)
}

func (m *mockDataSource) Commit(dctx DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

func (m *mockDataSource) Rollback(dctx DataContext) error {
	args := m.Called(dctx)
	return args.Error(0)
}

func Test_WithQuerier(t *testing.T) {
	expected, err := pgxmock.NewPool()
	require.NoError(t, err)
	ctx := WithQuerier(context.Background(), expected)
	querier, ok := ctx.Value(contextKey.querier).(Querier)
	assert.True(t, ok)
	assert.Equal(t, expected, querier)
}

func Test_QuerierFromCtx(t *testing.T) {
	expected, err := pgxmock.NewPool()
	require.NoError(t, err)
	ctx := context.WithValue(context.Background(), contextKey.querier, expected)
	querier, ok := QuerierFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, querier)
}

func Test_dbContext_Ctx(t *testing.T) {
	ctx := context.TODO()
	dctx := dbContext{Context: ctx}
	assert.Equal(t, ctx, dctx.Ctx())
}
