package user

import (
	"context"
	"errors"

	"strv-template-backend-go-api/database/sql"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/domain/user/postgres/query"
	"strv-template-backend-go-api/types"
	"strv-template-backend-go-api/types/id"

	"github.com/jackc/pgx/v5"
)

// Repository represents user data layer.
// Every model that is returned is converted to domain model using user factory.
type Repository struct {
	dataSource  sql.DataSource
	userFactory domuser.Factory
}

// NewRepository returns new instance of a user repository.
func NewRepository(dataSource sql.DataSource, userFactory domuser.Factory) (*Repository, error) {
	if dataSource == nil {
		return nil, errors.New("invalid data source")
	}
	return &Repository{
		dataSource:  dataSource,
		userFactory: userFactory,
	}, nil
}

// Create creates user in the repository.
func (r *Repository) Create(ctx context.Context, user *domuser.User) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		err := sql.Exec(dctx, query.Create, pgx.NamedArgs{
			"id":            user.ID,
			"referrer_id":   user.ReferrerID,
			"name":          user.Name,
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"role":          user.Role,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		})
		if err != nil {
			if sql.IsUnique(err, query.ConstraintUniqueUserEmail) {
				return domuser.ErrUserEmailAlreadyExists
			}
			if sql.IsForeignKey(err, query.ConstraintReferrerID) {
				return domuser.ErrReferrerNotFound
			}
			return err
		}
		return nil
	})
}

// Read reads the user from the repository.
func (r *Repository) Read(ctx context.Context, userID id.User) (*domuser.User, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domuser.User, error) {
		return r.read(dctx, userID)
	})
}

func (r *Repository) read(dctx sql.DataContext, userID id.User) (*domuser.User, error) {
	user, err := sql.ReadValue[user](dctx, query.Read, pgx.NamedArgs{
		"id": userID,
	})
	if err != nil {
		if sql.IsNotFound(err) {
			return nil, domuser.ErrUserNotFound
		}
		return nil, err
	}
	return user.ToUser(r.userFactory), nil
}

// ReadByEmail reads the user by email from the repository.
func (r *Repository) ReadByEmail(ctx context.Context, email types.Email) (*domuser.User, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domuser.User, error) {
		user, err := sql.ReadValue[user](dctx, query.ReadByEmail, pgx.NamedArgs{
			"email": email,
		})
		if err != nil {
			if sql.IsNotFound(err) {
				return nil, domuser.ErrUserNotFound
			}
			return nil, err
		}
		return user.ToUser(r.userFactory), nil
	})
}

// List lists users from the repository.
func (r *Repository) List(ctx context.Context) ([]domuser.User, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) ([]domuser.User, error) {
		dbUsers, err := sql.List[user](dctx, query.List)
		if err != nil {
			return nil, err
		}
		users := make([]domuser.User, 0, len(dbUsers))
		for _, u := range dbUsers {
			users = append(users, *u.ToUser(r.userFactory))
		}
		return users, nil
	})
}

// Update reads the user, calls external update function and updates the user in the repository.
func (r *Repository) Update(ctx context.Context, userID id.User, updateFn domuser.UpdateFunc) error {
	return sql.WithTransaction(ctx, r.dataSource, func(dctx sql.DataContext) error {
		originalUser, err := r.read(dctx, userID)
		if err != nil {
			return err
		}

		newUser, err := updateFn(originalUser)
		if err != nil {
			return err
		}

		return sql.ExecOne(dctx, query.Update, pgx.NamedArgs{
			"id":            newUser.ID,
			"referrer_id":   newUser.ReferrerID,
			"name":          newUser.Name,
			"email":         newUser.Email,
			"password_hash": newUser.PasswordHash,
			"role":          newUser.Role,
			"created_at":    newUser.CreatedAt,
			"updated_at":    newUser.UpdatedAt,
		})
	})
}
