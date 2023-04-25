package graph

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"strv-template-backend-go-api/api/graphql/graph/model"
	"strv-template-backend-go-api/database/sql"
	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	userdataloader "strv-template-backend-go-api/domain/user/postgres/dataloader"
	"strv-template-backend-go-api/domain/user/postgres/dataloader/query"
	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"
	utilctx "strv-template-backend-go-api/util/context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserName     = "Jozko Dlouhy"
	testUserEmail    = types.Email("jozko.dlouhy@gmail.com")
	testUserPassword = types.Password("Topsecret1")
	testRefreshToken = id.RefreshToken("5asd4a6d4a36d45as36da")
)

func newUser(t *testing.T) *domuser.User {
	t.Helper()
	factory, err := domuser.NewFactory(&mockHasher{}, &mockTimeSource{})
	require.NoError(t, err)
	now := time.Now()
	user := factory.NewUserFromFields(
		id.NewUser(),
		nil,
		testUserName,
		testUserEmail.String(),
		[]byte("sa65da263sa56as3"),
		string(domuser.RoleUser),
		now,
		now,
	)
	return user
}

func newSession(t *testing.T, claims domsession.Claims) *domsession.Session {
	t.Helper()
	mockTimeSource := &mockTimeSource{}
	mockTimeSource.On("Now").Return(time.Now())
	factory, err := domsession.NewFactory([]byte("abc123"), mockTimeSource, time.Hour, time.Hour)
	require.NoError(t, err)
	session, err := factory.NewSession(claims)
	require.NoError(t, err)
	return session
}

func Test_mutationResolver_CreateUser(t *testing.T) {
	input := model.CreateUserInput{
		Name:       testUserName,
		Email:      model.Email(testUserEmail),
		Password:   testUserPassword,
		ReferrerID: nil,
	}
	createUserInput := domuser.CreateUserInput{
		Name:       input.Name,
		Email:      types.Email(input.Email),
		Password:   input.Password,
		ReferrerID: nil,
	}
	var user *domuser.User
	var session *domsession.Session
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		input model.CreateUserInput
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.CreateUserResponse
		expectedErr error
	}{
		{
			name: "success",
			args: args{input: input},
			mocks: func() mocks {
				mocks := newMocks()
				user = newUser(t)
				session = newSession(t, domsession.Claims{
					UserID: id.NewUser(),
					Custom: domsession.CustomClaims{UserRole: domuser.RoleUser},
				})
				mocks.userService.On("Create", createUserInput).Return(user, session, nil)
				return mocks
			}(),
			expected: &model.CreateUserResponse{
				User:    model.FromUser(user),
				Session: model.FromSession(session),
			},
			expectedErr: nil,
		},
		{
			name: "failure:create-user",
			args: args{input: input},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Create", createUserInput).Return((*domuser.User)(nil), (*domsession.Session)(nil), serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
		{
			name:        "failure:new-create-user-input",
			args:        args{input: model.CreateUserInput{Email: model.Email(testUserEmail)}},
			mocks:       newMocks(),
			expected:    nil,
			expectedErr: apierrors.NewInvalidBodyError(domuser.ErrInvalidUserName, "new create user input").WithPublicMessage(domuser.ErrInvalidUserName.Error()),
		},
		{
			name:        "failure:new-email",
			args:        args{input: model.CreateUserInput{}},
			mocks:       newMocks(),
			expected:    nil,
			expectedErr: errors.New("new email: Key: '' Error:Field validation for '' failed on the 'email' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mutationResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.CreateUser(context.Background(), tt.args.input)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_mutationResolver_ChangeUserPassword(t *testing.T) {
	input := model.ChangeUserPasswordInput{
		NewPassword: testUserPassword,
		OldPassword: "Topsecret2",
	}
	userID := id.NewUser()
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		ctx   context.Context
		input model.ChangeUserPasswordInput
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.ChangeUserPasswordResponse
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx:   utilctx.WithUserID(context.Background(), userID),
				input: input,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ChangePassword", userID, input.OldPassword, input.NewPassword).Return(nil)
				return mocks
			}(),
			expected:    &model.ChangeUserPasswordResponse{Message: "password changed successfully"},
			expectedErr: nil,
		},
		{
			name: "failure:change-password",
			args: args{
				ctx:   utilctx.WithUserID(context.Background(), userID),
				input: input,
			},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ChangePassword", userID, input.OldPassword, input.NewPassword).Return(serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
		{
			name: "failure:missing-user-id",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			mocks:       newMocks(),
			expected:    nil,
			expectedErr: errMissingUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mutationResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.ChangeUserPassword(tt.args.ctx, tt.args.input)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_mutationResolver_CreateSession(t *testing.T) {
	var user *domuser.User
	var session *domsession.Session
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		input model.CreateSessionInput
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.CreateSessionResponse
		expectedErr error
	}{
		{
			name: "success",
			args: args{input: model.CreateSessionInput{
				Email:    model.Email(testUserEmail),
				Password: testUserPassword,
			}},
			mocks: func() mocks {
				mocks := newMocks()
				session = newSession(t, domsession.Claims{
					UserID: id.NewUser(),
					Custom: domsession.CustomClaims{UserRole: domuser.RoleUser},
				})
				user = newUser(t)
				mocks.sessionService.On("Create", testUserEmail, testUserPassword).Return(session, user, nil)
				return mocks
			}(),
			expected: &model.CreateSessionResponse{
				User:    model.FromUser(user),
				Session: model.FromSession(session),
			},
			expectedErr: nil,
		},
		{
			name: "failure:create",
			args: args{input: model.CreateSessionInput{
				Email:    model.Email(testUserEmail),
				Password: testUserPassword,
			}},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Create", testUserEmail, testUserPassword).Return((*domsession.Session)(nil), (*domuser.User)(nil), serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
		{
			name:        "failure:new-email",
			args:        args{input: model.CreateSessionInput{}},
			mocks:       newMocks(),
			expected:    nil,
			expectedErr: errors.New("new email: Key: '' Error:Field validation for '' failed on the 'email' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mutationResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.CreateSession(context.Background(), tt.args.input)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_mutationResolver_RefreshSession(t *testing.T) {
	var session *domsession.Session
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		input model.RefreshSessionInput
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.RefreshSessionResponse
		expectedErr error
	}{
		{
			name: "success",
			args: args{input: model.RefreshSessionInput{RefreshToken: string(testRefreshToken)}},
			mocks: func() mocks {
				mocks := newMocks()
				session = newSession(t, domsession.Claims{
					UserID: id.NewUser(),
					Custom: domsession.CustomClaims{UserRole: domuser.RoleUser},
				})
				mocks.sessionService.On("Refresh", testRefreshToken).Return(session, nil)
				return mocks
			}(),
			expected:    model.FromSession(session),
			expectedErr: nil,
		},
		{
			name: "failure:refresh",
			args: args{input: model.RefreshSessionInput{RefreshToken: string(testRefreshToken)}},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Refresh", testRefreshToken).Return((*domsession.Session)(nil), serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mutationResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.RefreshSession(context.Background(), tt.args.input)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_mutationResolver_DestroySession(t *testing.T) {
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		input model.DestroySessionInput
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.DestroySessionResponse
		expectedErr error
	}{
		{
			name: "success",
			args: args{input: model.DestroySessionInput{RefreshToken: string(testRefreshToken)}},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Destroy", testRefreshToken).Return(nil)
				return mocks
			}(),
			expected:    &model.DestroySessionResponse{Message: "session destroyed"},
			expectedErr: nil,
		},
		{
			name: "failure:refresh",
			args: args{input: model.DestroySessionInput{RefreshToken: string(testRefreshToken)}},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Destroy", testRefreshToken).Return(serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mutationResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.DestroySession(context.Background(), tt.args.input)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_queryResolver_UserMe(t *testing.T) {
	userID := id.NewUser()
	var user *domuser.User
	serviceErr := apierrors.NewUnknownError(errTest, "")

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		mocks       mocks
		expected    *model.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{ctx: utilctx.WithUserID(context.Background(), userID)},
			mocks: func() mocks {
				mocks := newMocks()
				user = newUser(t)
				mocks.userService.On("Read", userID).Return(user, nil)
				return mocks
			}(),
			expected:    model.FromUser(user),
			expectedErr: nil,
		},
		{
			name: "failure:read",
			args: args{ctx: utilctx.WithUserID(context.Background(), userID)},
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Read", userID).Return((*domuser.User)(nil), serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
		{
			name:        "failure:missing-user-id",
			args:        args{ctx: context.Background()},
			mocks:       newMocks(),
			expected:    nil,
			expectedErr: errMissingUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &queryResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.UserMe(tt.args.ctx)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_queryResolver_Users(t *testing.T) {
	var users []domuser.User
	serviceErr := apierrors.NewUnknownError(errTest, "")

	tests := []struct {
		name        string
		mocks       mocks
		expected    []model.User
		expectedErr error
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				users = []domuser.User{*newUser(t), *newUser(t)}
				mocks.userService.On("List").Return(users, nil)
				return mocks
			}(),
			expected:    model.FromUsers(users),
			expectedErr: nil,
		},
		{
			name: "failure:list",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("List").Return([]domuser.User(nil), serviceErr)
				return mocks
			}(),
			expected:    nil,
			expectedErr: serviceErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &queryResolver{&Resolver{
				userService:    tt.mocks.userService,
				sessionService: tt.mocks.sessionService,
			}}
			result, err := resolver.Users(context.Background())
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_userResolver_Referrer(t *testing.T) {
	referrer := newUser(t)
	referrerID := uuid.MustParse(referrer.ID.String())

	type args struct {
		ctx context.Context
		obj *model.User
	}
	tests := []struct {
		name        string
		args        args
		expected    *model.User
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx: func() context.Context {
					querier, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(queryMatcher))
					require.NoError(t, err)
					dataSource := &mockDataSource{}
					dctx := dbContext{Context: sql.WithQuerier(context.Background(), querier)}
					dataSource.On("AcquireConnCtx").Return(dctx, nil)
					dataSource.On("ReleaseConnCtx", dctx).Return(nil)
					queryArgs := pgx.NamedArgs{
						"ids": []string{referrerID.String()},
					}
					rows := pgxmock.NewRows([]string{"id", "referrer_id", "name", "email", "password_hash", "role", "created_at", "updated_at"})
					rows = rows.AddRow(referrer.ID, referrer.ReferrerID, referrer.Name, referrer.Email, referrer.PasswordHash, referrer.Role, referrer.CreatedAt, referrer.UpdatedAt)
					querier.ExpectQuery(query.ListUsersByIDs).WithArgs(queryArgs).WillReturnRows(rows)
					userLoader := userdataloader.New(dataSource)
					return userdataloader.WithLoader(context.Background(), userLoader)
				}(),
				obj: &model.User{ReferrerID: &referrerID},
			},
			expected:    model.FromUser(referrer),
			expectedErr: nil,
		},
		{
			name: "failure:read-user",
			args: args{
				ctx: func() context.Context {
					dataSource := &mockDataSource{}
					dataSource.On("AcquireConnCtx").Return(dbContext{}, errTest)
					userLoader := userdataloader.New(dataSource)
					return userdataloader.WithLoader(context.Background(), userLoader)
				}(),
				obj: &model.User{ReferrerID: &referrerID},
			},
			expected:    nil,
			expectedErr: fmt.Errorf("reading user: %w", fmt.Errorf("acquiring connection: %w", errTest)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &userResolver{}
			result, err := resolver.Referrer(tt.args.ctx, tt.args.obj)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
