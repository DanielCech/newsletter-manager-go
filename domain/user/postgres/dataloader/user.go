package dataloader

import (
	"context"
	"errors"

	"strv-template-backend-go-api/database/sql"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/domain/user/postgres/dataloader/query"
	apierrors "strv-template-backend-go-api/types/errors"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v5"
)

// UserReader defines methods for fetching of users.
type UserReader struct {
	dataSource sql.DataSource
}

// NewUserReader returns new instance of UserReader.
func NewUserReader(dataSource sql.DataSource) *UserReader {
	return &UserReader{dataSource: dataSource}
}

// ListUsersByIDs returns all users that suit the given key.
// Keys represent user IDs.
func (u *UserReader) ListUsersByIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	output := make([]*dataloader.Result, len(keys))
	userIDs := make([]string, len(keys))
	for i, key := range keys {
		userIDs[i] = key.String()
	}

	users, err := sql.WithConnectionResult[[]domuser.User](ctx, u.dataSource, func(dctx sql.DataContext) ([]domuser.User, error) {
		users, err := sql.List[domuser.User](dctx, query.ListUsersByIDs, pgx.NamedArgs{
			"ids": userIDs,
		})
		if err != nil {
			return nil, apierrors.NewUnknownError(err, "listing users by ids")
		}
		return users, nil
	})
	if err != nil {
		for i := range userIDs {
			output[i] = &dataloader.Result{
				Data:  nil,
				Error: err,
			}
		}
		return output
	}

	userByID := map[string]*domuser.User{}
	for i := range users {
		userByID[users[i].ID.String()] = &users[i]
	}

	for i, key := range keys {
		user, ok := userByID[key.String()]
		if !ok {
			err = errors.New("user not found")
			output[i] = &dataloader.Result{
				Data: nil,
				Error: apierrors.NewNotFoundError(err, "").WithPublicMessage(err.Error()).WithData(map[string]any{
					"id": key.String(),
				}),
			}
			continue
		}
		output[i] = &dataloader.Result{Data: user, Error: nil}
	}

	return output
}
