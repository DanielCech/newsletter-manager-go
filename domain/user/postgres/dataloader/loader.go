package dataloader

import (
	"context"
	"errors"
	"fmt"

	"strv-template-backend-go-api/database/sql"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"

	"github.com/graph-gophers/dataloader"
)

var (
	contextKey = struct {
		loader ctxKeyLoader
	}{}
)

type (
	ctxKeyLoader struct{}
)

// Loader contains data loaders which are useful for fetching of users.
type Loader struct {
	User *dataloader.Loader
}

// New returns new instance of user Loader.
func New(dataSource sql.DataSource) *Loader {
	userReader := NewUserReader(dataSource)
	return &Loader{
		User: dataloader.NewBatchedLoader(userReader.ListUsersByIDs),
	}
}

// WithLoader passes loader to the context.
func WithLoader(ctx context.Context, loader *Loader) context.Context {
	return context.WithValue(ctx, contextKey.loader, loader)
}

// LoaderFromCtx gets loader from the context.
func LoaderFromCtx(ctx context.Context) (*Loader, bool) {
	loader, ok := ctx.Value(contextKey.loader).(*Loader)
	return loader, ok
}

// ReadUser reads user by use of data loader.
func ReadUser(ctx context.Context, userID id.User) (*domuser.User, error) {
	loader, ok := LoaderFromCtx(ctx)
	if !ok {
		return nil, errors.New("missing data loader")
	}
	thunk := loader.User.Load(ctx, dataloader.StringKey(userID.String()))

	loadedUser, err := thunk()
	if err != nil {
		return nil, err
	}

	result, ok := loadedUser.(*domuser.User)
	if !ok {
		return nil, fmt.Errorf("wrong type returned from loader: %T", loadedUser)
	}

	return result, nil
}
