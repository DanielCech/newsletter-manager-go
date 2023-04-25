package dataloader

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"strv-template-backend-go-api/database/sql"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/domain/user/postgres/dataloader/query"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func newUser() *domuser.User {
	now := time.Now()
	referrerID := id.NewUser()
	return &domuser.User{
		ID:           id.NewUser(),
		ReferrerID:   &referrerID,
		Name:         "Jozko Dlouhy",
		Email:        "jozko.dlouhy@gmail.com",
		PasswordHash: []byte("45s4das545dadsa25"),
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func Test_NewUserReader(t *testing.T) {
	userReader := NewUserReader(&mockDataSource{})
	assert.NotEmpty(t, userReader)
}

func Test_UserReader_ListUsersByIDs(t *testing.T) {
	ctx := context.Background()
	user := newUser()
	unknownUserID := id.NewUser()

	type args struct {
		keys dataloader.Keys
	}
	tests := []struct {
		name     string
		args     args
		mocks    mocks
		expected []*dataloader.Result
	}{
		{
			name: "success",
			args: args{keys: dataloader.NewKeysFromStrings([]string{user.ID.String()})},
			mocks: func() mocks {
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
				return mocks
			}(),
			expected: []*dataloader.Result{
				{
					Data:  user,
					Error: nil,
				},
			},
		},
		{
			name: "success-partial",
			args: args{keys: dataloader.NewKeysFromStrings([]string{user.ID.String(), unknownUserID.String()})},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"ids": []string{user.ID.String(), unknownUserID.String()},
				}
				rows := pgxmock.NewRows([]string{"id", "referrer_id", "name", "email", "password_hash", "role", "created_at", "updated_at"})
				rows = rows.AddRow(user.ID, user.ReferrerID, user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
				mocks.querier.ExpectQuery(query.ListUsersByIDs).WithArgs(queryArgs).WillReturnRows(rows)
				return mocks
			}(),
			expected: []*dataloader.Result{
				{
					Data:  user,
					Error: nil,
				},
				{
					Data: nil,
					Error: func() error {
						err := errors.New("user not found")
						return apierrors.NewNotFoundError(err, "").WithPublicMessage(err.Error()).WithData(map[string]any{
							"id": unknownUserID.String(),
						})
					}(),
				},
			},
		},
		{
			name: "failure:listing-users",
			args: args{keys: dataloader.NewKeysFromStrings([]string{user.ID.String()})},
			mocks: func() mocks {
				mocks := newMocks(t)
				dctx := dbContext{Context: sql.WithQuerier(ctx, mocks.querier)}
				mocks.dataSource.On("AcquireConnCtx").Return(dctx, nil)
				mocks.dataSource.On("ReleaseConnCtx", dctx).Return(nil)
				queryArgs := pgx.NamedArgs{
					"ids": []string{user.ID.String()},
				}
				mocks.querier.ExpectQuery(query.ListUsersByIDs).WithArgs(queryArgs).WillReturnError(errTest)
				return mocks
			}(),
			expected: []*dataloader.Result{
				{
					Data:  nil,
					Error: apierrors.NewUnknownError(fmt.Errorf("scany: query multiple result rows: %w", errTest), "listing users by ids"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userReader := NewUserReader(tt.mocks.dataSource)
			result := userReader.ListUsersByIDs(ctx, tt.args.keys)
			assert.Equal(t, tt.expected, result)
			tt.mocks.assertExpectations(t)
		})
	}
}
