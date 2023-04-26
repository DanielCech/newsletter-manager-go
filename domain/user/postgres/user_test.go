package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"newsletter-manager-go/database/sql"
	domuser "newsletter-manager-go/domain/user"
	"newsletter-manager-go/domain/user/postgres/query"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
	"newsletter-manager-go/util/timesource"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("test error")
)

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
	hasher     *mockHasher
	timeSource *mockTimeSource
	dataSource *mockDataSource
	querier    pgxmock.PgxPoolIface
}

func newMocks(t *testing.T) mocks {
	t.Helper()
	querier, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(queryMatcher))
	require.NoError(t, err)
	return mocks{
		hasher:     &mockHasher{},
		timeSource: &mockTimeSource{},
		dataSource: &mockDataSource{},
		querier:    querier,
	}
}

func (m mocks) assertExpectations(t *testing.T) {
	t.Helper()
	m.hasher.AssertExpectations(t)
	m.timeSource.AssertExpectations(t)
	m.dataSource.AssertExpectations(t)
	assert.NoError(t, m.querier.ExpectationsWereMet())
}

func newFactory(t *testing.T, hasher domuser.Hasher, timeSource timesource.TimeSource) domuser.Factory {
	t.Helper()
	userFactory, err := domuser.NewFactory(
		hasher,
		timeSource,
	)
	require.NoError(t, err)
	return userFactory
}

func newRepository(t *testing.T, mocks mocks) *Repository {
	t.Helper()
	userFactory := newFactory(t, mocks.hasher, mocks.timeSource)
	repository, err := NewRepository(mocks.dataSource, userFactory)
	require.NoError(t, err)
	return repository
}

func newUser(t *testing.T) *domuser.User {
	t.Helper()
	factory := newFactory(t, &mockHasher{}, &mockTimeSource{})
	now := time.Now()
	referrerID := id.NewUser()
	return factory.NewUserFromFields(
		id.NewUser(),
		&referrerID,
		"Jozko Dlouhy",
		"jozko.dlouhy@gmail.com",
		[]byte("45s4das545dadsa25"),
		"user",
		now,
		now,
	)
}

func newUserRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"referrer_id",
		"name",
		"email",
		"password_hash",
		"role",
		"created_at",
		"updated_at",
	})
}

func addUserToRows(user *domuser.User, rows *pgxmock.Rows) *pgxmock.Rows {
	return rows.AddRow(
		user.ID,
		user.ReferrerID,
		user.Name,
		string(user.Email),
		user.PasswordHash,
		string(user.Role),
		user.CreatedAt,
		user.UpdatedAt,
	)
}

func userQueryArgs(user *domuser.User) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":            user.ID,
		"referrer_id":   user.ReferrerID,
		"name":          user.Name,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"role":          user.Role,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}
}

func Test_NewRepository(t *testing.T) {
	t.Helper()
	mocks := newMocks(t)
	userFactory := newFactory(t, mocks.hasher, mocks.timeSource)
	repository, err := NewRepository(mocks.dataSource, userFactory)
	require.NoError(t, err)
	assert.Equal(t, mocks.dataSource, repository.dataSource)
	assert.Equal(t, userFactory, repository.userFactory)

	repository, err = NewRepository(nil, userFactory)
	assert.EqualError(t, err, "invalid data source")
	assert.Empty(t, repository)
}

func Test_Repository_Create(t *testing.T) {
	ctx := context.Background()
	user := newUser(t)

	type args struct {
		user *domuser.User
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expectedErr error
	}{
		{
			name: "success",
			args: args{user: user},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				result := pgxmock.NewResult("CREATE", 1)
				mocks.querier.ExpectExec(query.Create).WithArgs(userQueryArgs(user)).WillReturnResult(result)
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:create-user",
			args: args{user: user},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectExec(query.Create).WithArgs(userQueryArgs(user)).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:create-user/foreign-key",
			args: args{user: user},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				err := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation, ConstraintName: query.ConstraintReferrerID}
				mocks.querier.ExpectExec(query.Create).WithArgs(userQueryArgs(user)).WillReturnError(err)
				return mocks
			}(),
			expectedErr: domuser.ErrReferrerNotFound,
		},
		{
			name: "failure:create-user/unique",
			args: args{user: user},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				err := &pgconn.PgError{Code: pgerrcode.UniqueViolation, ConstraintName: query.ConstraintUniqueUserEmail}
				mocks.querier.ExpectExec(query.Create).WithArgs(userQueryArgs(user)).WillReturnError(err)
				return mocks
			}(),
			expectedErr: domuser.ErrUserEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.Create(ctx, tt.args.user)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_Read(t *testing.T) {
	ctx := context.Background()
	userID := id.NewUser()
	user := newUser(t)
	queryArgs := pgx.NamedArgs{
		"id": userID,
	}

	type args struct {
		userID id.User
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domuser.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{userID: userID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				rows := newUserRows()
				rows = addUserToRows(user, rows)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgs).WillReturnRows(rows)
				return mocks
			}(),
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "failure:read-user",
			args: args{userID: userID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgs).WillReturnError(errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("scany: query one result row: %w", errTest),
		},
		{
			name: "failure:read-user/not-found",
			args: args{userID: userID},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgs).WillReturnError(pgx.ErrNoRows)
				return mocks
			}(),
			expected:    nil,
			expectedErr: domuser.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			user, err := repository.Read(ctx, tt.args.userID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_ReadByEmail(t *testing.T) {
	ctx := context.Background()
	email := types.Email("jozko.dlouhy@gmail.com")
	user := newUser(t)
	queryArgs := pgx.NamedArgs{
		"email": email,
	}

	type args struct {
		email types.Email
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *domuser.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{email: email},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				rows := newUserRows()
				rows = addUserToRows(user, rows)
				mocks.querier.ExpectQuery(query.ReadByEmail).WithArgs(queryArgs).WillReturnRows(rows)
				return mocks
			}(),
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "failure:read-user",
			args: args{email: email},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectQuery(query.ReadByEmail).WithArgs(queryArgs).WillReturnError(errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("scany: query one result row: %w", errTest),
		},
		{
			name: "failure:read-user/not-found",
			args: args{email: email},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectQuery(query.ReadByEmail).WithArgs(queryArgs).WillReturnError(pgx.ErrNoRows)
				return mocks
			}(),
			expected:    nil,
			expectedErr: domuser.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			user, err := repository.ReadByEmail(ctx, tt.args.email)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_List(t *testing.T) {
	ctx := context.Background()
	users := []domuser.User{*newUser(t), *newUser(t)}

	tests := []struct {
		name        string
		mocks       mocks
		expected    []domuser.User
		expectedErr error
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				rows := newUserRows()
				rows = addUserToRows(&users[0], rows)
				rows = addUserToRows(&users[1], rows)
				mocks.querier.ExpectQuery(query.List).WillReturnRows(rows)
				return mocks
			}(),
			expected:    users,
			expectedErr: nil,
		},
		{
			name: "failure:list-users",
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				mocks.querier.ExpectQuery(query.List).WillReturnError(errTest)
				return mocks
			}(),
			expected:    nil,
			expectedErr: fmt.Errorf("scany: query multiple result rows: %w", errTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			users, err := repository.List(ctx)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, users)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Repository_Update(t *testing.T) {
	ctx := context.Background()
	user := newUser(t)
	queryArgsID := pgx.NamedArgs{
		"id": user.ID,
	}

	type args struct {
		id       id.User
		updateFn domuser.UpdateFunc
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
				id: user.ID,
				updateFn: func(u *domuser.User) (*domuser.User, error) {
					return u, nil
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Commit", dctx).Return(nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read user.
				rows := newUserRows()
				rows = addUserToRows(user, rows)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgsID).WillReturnRows(rows)
				// Update user.
				mocks.querier.ExpectExec(query.Update).WithArgs(userQueryArgs(user)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				return mocks
			}(),
			expectedErr: nil,
		},
		{
			name: "failure:update-user",
			args: args{
				id: user.ID,
				updateFn: func(u *domuser.User) (*domuser.User, error) {
					return u, nil
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read user.
				rows := newUserRows()
				rows = addUserToRows(user, rows)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgsID).WillReturnRows(rows)
				// Update user.
				mocks.querier.ExpectExec(query.Update).WithArgs(userQueryArgs(user)).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:update-user",
			args: args{
				id: user.ID,
				updateFn: func(u *domuser.User) (*domuser.User, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read user.
				rows := newUserRows()
				rows = addUserToRows(user, rows)
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgsID).WillReturnRows(rows)
				return mocks
			}(),
			expectedErr: errTest,
		},
		{
			name: "failure:update-user",
			args: args{
				id: user.ID,
				updateFn: func(u *domuser.User) (*domuser.User, error) {
					return nil, errTest
				},
			},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("Begin").Return(dctx, nil)
				mocks.dataSource.On("Rollback", dctx).Return(nil)
				// Read user.
				mocks.querier.ExpectQuery(query.Read).WithArgs(queryArgsID).WillReturnError(errTest)
				return mocks
			}(),
			expectedErr: fmt.Errorf("scany: query one result row: %w", errTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository(t, tt.mocks)
			err := repository.Update(ctx, tt.args.id, tt.args.updateFn)
			assert.Equal(t, tt.expectedErr, err)
			tt.mocks.assertExpectations(t)
		})
	}
}
