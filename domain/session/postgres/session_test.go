package session

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"strv-template-backend-go-api/database/sql"
	domsession "strv-template-backend-go-api/domain/session"
	"strv-template-backend-go-api/domain/session/postgres/query"
	"strv-template-backend-go-api/types/id"
	"strv-template-backend-go-api/util/timesource"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

type mockTimeSource struct {
	mock.Mock
}

func (m *mockTimeSource) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mockDataSource struct {
	mock.Mock
}

type dbContext struct {
	context.Context
}

func (d dbContext) Ctx() context.Context {
	return d.Context
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

var queryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	if strings.Compare(expectedSQL, actualSQL) == 0 {
		return nil
	}
	return fmt.Errorf(`could not match actual sql: "%s" with expected: "%s"`, actualSQL, expectedSQL)
})

type mocks struct {
	timeSource *mockTimeSource
	dataSource *mockDataSource
	querier    pgxmock.PgxPoolIface
}

func newMocks(t *testing.T) mocks {
	t.Helper()
	querier, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(queryMatcher))
	require.NoError(t, err)
	return mocks{
		timeSource: &mockTimeSource{},
		dataSource: &mockDataSource{},
		querier:    querier,
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.timeSource.AssertExpectations(t)
	m.dataSource.AssertExpectations(t)
	assert.NoError(t, m.querier.ExpectationsWereMet())
}

func newFactory(t *testing.T, timeSource timesource.TimeSource) domsession.Factory {
	t.Helper()
	sessionFactory, err := domsession.NewFactory(
		[]byte("abc123"),
		timeSource,
		time.Hour,
		time.Hour,
	)
	require.NoError(t, err)
	return sessionFactory
}

func newRepository(t *testing.T, mocks mocks) *Repository {
	t.Helper()
	sessionFactory := newFactory(t, mocks.timeSource)
	repository, err := NewRepository(mocks.dataSource, sessionFactory)
	require.NoError(t, err)
	return repository
}

func newRefreshTokenRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"user_id",
		"expires_at",
		"created_at",
	})
}

func addRefreshTokenToRows(refreshToken *domsession.RefreshToken, rows *pgxmock.Rows) *pgxmock.Rows {
	return rows.AddRow(
		refreshToken.ID,
		refreshToken.AuthorID,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
	)
}

func refreshTokenQueryArgs(refreshToken *domsession.RefreshToken) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":         refreshToken.ID,
		"user_id":    refreshToken.AuthorID,
		"expires_at": refreshToken.ExpiresAt,
		"created_at": refreshToken.CreatedAt,
	}
}

func Test_NewRepository(t *testing.T) {
	t.Helper()
	mocks := newMocks(t)
	sessionFactory := newFactory(t, mocks.timeSource)
	repository, err := NewRepository(mocks.dataSource, sessionFactory)
	require.NoError(t, err)
	assert.Equal(t, mocks.dataSource, repository.dataSource)
	assert.Equal(t, sessionFactory, repository.sessionFactory)

	repository, err = NewRepository(nil, sessionFactory)
	assert.EqualError(t, err, "invalid data source")
	assert.Empty(t, repository)
}

func Test_Repository_CreateRefreshToken(t *testing.T) {
	ctx := context.Background()
	factory := newFactory(t, timesource.DefaultTimeSource{})
	refreshToken, err := factory.NewRefreshToken(id.NewUser())
	require.NoError(t, err)

	type args struct {
		refreshToken *domsession.RefreshToken
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{refreshToken: refreshToken},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				result := pgxmock.NewResult("CREATE", 1)
				mocks.querier.ExpectExec(query.CreateRefreshToken).WithArgs(refreshTokenQueryArgs(refreshToken)).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:create-refresh-token",
			args: args{refreshToken: refreshToken},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectExec(query.CreateRefreshToken).WithArgs(refreshTokenQueryArgs(refreshToken)).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.CreateRefreshToken(ctx, tt.args.refreshToken)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_Refresh(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	refreshTokenID := id.RefreshToken("5asd4a6d4a36d45as36da")
	refreshToken := &domsession.RefreshToken{
		ID:        refreshTokenID,
		AuthorID:  id.NewUser(),
		ExpiresAt: now,
		CreatedAt: now,
	}
	queryArgsID := pgx.NamedArgs{
		"id": refreshTokenID,
	}

	type args struct {
		refreshTokenID id.RefreshToken
		refreshFn      domsession.RefreshFunc
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return token, nil
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Commit", dctx).Return(nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				rows := newRefreshTokenRows()
				rows = addRefreshTokenToRows(refreshToken, rows)
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnRows(rows)
				// Delete old refresh token.
				result := pgxmock.NewResult("DELETE", 1)
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgsID).WillReturnResult(result)
				// Create new refresh token.
				result = pgxmock.NewResult("CREATE", 1)
				mocks.querier.ExpectExec(query.CreateRefreshToken).WithArgs(refreshTokenQueryArgs(refreshToken)).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:create-refresh-token",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return token, nil
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				rows := newRefreshTokenRows()
				rows = addRefreshTokenToRows(refreshToken, rows)
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnRows(rows)
				// Delete old refresh token.
				result := pgxmock.NewResult("DELETE", 1)
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgsID).WillReturnResult(result)
				// Create new refresh token.
				mocks.querier.ExpectExec(query.CreateRefreshToken).WithArgs(refreshTokenQueryArgs(refreshToken)).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:refresh-fn",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				rows := newRefreshTokenRows()
				rows = addRefreshTokenToRows(refreshToken, rows)
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnRows(rows)
				// Delete old refresh token.
				result := pgxmock.NewResult("DELETE", 1)
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgsID).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:delete-refresh-token",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				rows := newRefreshTokenRows()
				rows = addRefreshTokenToRows(refreshToken, rows)
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnRows(rows)
				// Delete old refresh token.
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgsID).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:read-refresh-token",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("scany: query one result row: %w", errTest),
		},
		{
			name: "failure:read-refresh-token/not-found",
			args: args{
				refreshTokenID: refreshTokenID,
				refreshFn: func(token *domsession.RefreshToken) (*domsession.RefreshToken, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read refresh token.
				mocks.querier.ExpectQuery(query.ReadRefreshToken).WithArgs(queryArgsID).WillReturnError(pgx.ErrNoRows)
				return mocks
			}(),
			expectedErr: domsession.ErrRefreshTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.Refresh(ctx, tt.args.refreshTokenID, tt.args.refreshFn)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_DeleteRefreshToken(t *testing.T) {
	ctx := context.Background()
	refreshTokenID := id.RefreshToken("5asd4a6d4a36d45as36da")

	type args struct {
		refreshTokenID id.RefreshToken
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"id": refreshTokenID,
				}
				result := pgxmock.NewResult("DELETE", 1)
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgs).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:delete-refresh-token",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"id": refreshTokenID,
				}
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgs).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:delete-refresh-token/not-found",
			args: args{refreshTokenID: refreshTokenID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"id": refreshTokenID,
				}
				mocks.querier.ExpectExec(query.DeleteRefreshToken).WithArgs(queryArgs).WillReturnError(pgx.ErrNoRows)
				return mocks
			}(),
			expectedErr: domsession.ErrRefreshTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.DeleteRefreshToken(ctx, tt.args.refreshTokenID)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_DeleteRefreshTokensByAuthorID(t *testing.T) {
	ctx := context.Background()
	AuthorID := id.NewUser()

	type args struct {
		AuthorID id.Author
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{AuthorID: AuthorID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"user_id": AuthorID,
				}
				result := pgxmock.NewResult("DELETE", 4)
				mocks.querier.ExpectExec(query.DeleteRefreshTokensByAuthorID).WithArgs(queryArgs).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:delete-refresh-tokens-by-user-id",
			args: args{AuthorID: AuthorID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"user_id": AuthorID,
				}
				mocks.querier.ExpectExec(query.DeleteRefreshTokensByAuthorID).WithArgs(queryArgs).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.DeleteRefreshTokensByAuthorID(ctx, tt.args.AuthorID)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}
