package sql

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

type mockPGXPool struct {
	mock.Mock
}

func (m *mockPGXPool) Acquire(_ context.Context) (*pgxpool.Conn, error) {
	args := m.Called()
	return args.Get(0).(*pgxpool.Conn), args.Error(1)
}

func (m *mockPGXPool) Begin(_ context.Context) (pgx.Tx, error) {
	args := m.Called()
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *mockPGXPool) Close() {
	m.Called()
}

func Test_Secret_SecretARN(t *testing.T) {
	arn := "64as6da4da6"
	secret := &Secret{ARN: &arn}
	assert.Equal(t, &arn, secret.SecretARN())
}

func Test_Secret_LocalPath(t *testing.T) {
	path := "./sa486da/asd42a"
	secret := &Secret{Path: &path}
	assert.Equal(t, &path, secret.LocalPath())
}

func Test_DSNValues_ConnString(t *testing.T) {
	dsnValues := DSNValues{
		Host:     "localhost",
		Port:     "5432",
		UserName: "root",
		Password: "Tosecret1",
		DBName:   "template",
		SSLMode:  "",
	}
	connString, err := dsnValues.ConnString()
	assert.NoError(t, err)
	assert.Equal(t, "postgres://root:Tosecret1@localhost:5432/template?sslmode=disable", connString)

	dsnValues.Port = "invalid"
	connString, err = dsnValues.ConnString()
	assert.EqualError(t, err, `parsing port: strconv.ParseInt: parsing "invalid": invalid syntax`)
	assert.Empty(t, connString)
}

func Test_Database_Close(t *testing.T) {
	mockPool := &mockPGXPool{}
	mockPool.On("Close").Return()
	database := Database{pool: mockPool}
	database.Close()
}

func Test_Database_AcquireConnCtx(t *testing.T) {
	conn := &pgxpool.Conn{}
	mockPool := &mockPGXPool{}
	mockPool.On("Acquire").Return(conn, nil).Once()
	database := Database{pool: mockPool}
	dctx, err := database.AcquireConnCtx(context.TODO())
	assert.NoError(t, err)
	querier, ok := dctx.Ctx().Value(contextKey.querier).(Querier)
	assert.True(t, ok)
	assert.Equal(t, conn, querier)

	mockPool.On("Acquire").Return((*pgxpool.Conn)(nil), errTest).Once()
	dctx, err = database.AcquireConnCtx(context.TODO())
	assert.Equal(t, errTest, err)
	assert.Empty(t, dctx)

	mockPool.AssertExpectations(t)
}

func Test_Database_ReleaseConnCtx(t *testing.T) {
	database := Database{pool: &mockPGXPool{}}
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, &pgxpool.Conn{})}

	err := database.ReleaseConnCtx(dctx)
	assert.NoError(t, err)

	dctx = dbContext{Context: context.TODO()}
	err = database.ReleaseConnCtx(dctx)
	assert.Equal(t, ErrNoConn, err)
}

func Test_Database_Begin(t *testing.T) {
	tx := &pgxpool.Tx{}
	mockPool := &mockPGXPool{}
	mockPool.On("Begin").Return(tx, nil).Once()
	database := Database{pool: mockPool}
	dctx, err := database.Begin(context.TODO())
	assert.NoError(t, err)
	querier, ok := dctx.Ctx().Value(contextKey.querier).(Querier)
	assert.True(t, ok)
	assert.Equal(t, tx, querier)

	mockPool.On("Begin").Return((*pgxpool.Tx)(nil), errTest).Once()
	dctx, err = database.Begin(context.TODO())
	assert.Equal(t, errTest, err)
	assert.Empty(t, dctx)

	mockPool.AssertExpectations(t)
}

func Test_Database_Commit(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	mockPool.ExpectCommit()
	database := Database{pool: mockPool}
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Commit(dctx)
	assert.NoError(t, err)

	mockPool.ExpectCommit().WillReturnError(pgx.ErrTxClosed)
	dctx = dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Commit(dctx)
	assert.NoError(t, err)

	mockPool.ExpectCommit().WillReturnError(errTest)
	dctx = dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Commit(dctx)
	assert.Equal(t, errTest, err)

	dctx = dbContext{Context: context.TODO()}
	err = database.Commit(dctx)
	assert.Equal(t, ErrNoTx, err)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_Database_Rollback(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	mockPool.ExpectRollback()
	database := Database{pool: mockPool}
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Rollback(dctx)
	assert.NoError(t, err)

	mockPool.ExpectRollback().WillReturnError(pgx.ErrTxClosed)
	dctx = dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Rollback(dctx)
	assert.NoError(t, err)

	mockPool.ExpectRollback().WillReturnError(errTest)
	dctx = dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, mockPool)}
	err = database.Rollback(dctx)
	assert.Equal(t, errTest, err)

	dctx = dbContext{Context: context.TODO()}
	err = database.Rollback(dctx)
	assert.Equal(t, ErrNoTx, err)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func Test_Connection(t *testing.T) {
	conn := &pgxpool.Conn{}
	dctx := dbContext{Context: context.WithValue(context.TODO(), contextKey.querier, conn)}
	q := querier(dctx)
	assert.Equal(t, conn, q)

	assert.PanicsWithError(t, ErrNoQuerier.Error(), func() {
		dctx = dbContext{Context: context.TODO()}
		_ = querier(dctx)
	})
}
