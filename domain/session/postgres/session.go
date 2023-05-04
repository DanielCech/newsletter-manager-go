package session

import (
	"context"
	"errors"

	"newsletter-manager-go/database/sql"
	domsession "newsletter-manager-go/domain/session"
	"newsletter-manager-go/domain/session/postgres/query"
	"newsletter-manager-go/types/id"

	"github.com/jackc/pgx/v5"
)

// Repository represents session data layer.
// Every model that is returned is converted to domain model using session factory.
type Repository struct {
	dataSource     sql.DataSource
	sessionFactory domsession.Factory
}

// NewRepository returns new instance of a session repository.
func NewRepository(dataSource sql.DataSource, sessionFactory domsession.Factory) (*Repository, error) {
	if dataSource == nil {
		return nil, errors.New("invalid data source")
	}
	return &Repository{
		dataSource:     dataSource,
		sessionFactory: sessionFactory,
	}, nil
}

// CreateRefreshToken creates refresh token in the repository.
func (r *Repository) CreateRefreshToken(ctx context.Context, refreshToken *domsession.RefreshToken) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		return r.createRefreshToken(dctx, refreshToken)
	})
}

func (r *Repository) createRefreshToken(dctx sql.DataContext, refreshToken *domsession.RefreshToken) error {
	return sql.ExecOne(dctx, query.CreateRefreshToken, pgx.NamedArgs{
		"id":         refreshToken.ID,
		"user_id":    refreshToken.AuthorID,
		"expires_at": refreshToken.ExpiresAt,
		"created_at": refreshToken.CreatedAt,
	})
}

// Refresh reads the refresh token and deletes it from the repository. Refresh token is converted into the domain model
// and is used as an argument for external refresh function. New refresh token is then created in the repository.
func (r *Repository) Refresh(ctx context.Context, refreshTokenID id.RefreshToken, refreshFn domsession.RefreshFunc) error {
	return sql.WithTransaction(ctx, r.dataSource, func(dctx sql.DataContext) error {
		token, err := sql.ReadValue[refreshToken](dctx, query.ReadRefreshToken, pgx.NamedArgs{
			"id": refreshTokenID,
		})
		if err != nil {
			if sql.IsNotFound(err) {
				return domsession.ErrRefreshTokenNotFound
			}
			return err
		}

		if err = r.deleteRefreshToken(dctx, token.ID); err != nil {
			return err
		}

		oldRefreshToken := token.ToRefreshToken(r.sessionFactory)
		newRefreshToken, err := refreshFn(oldRefreshToken)
		if err != nil {
			return err
		}

		if err = r.createRefreshToken(dctx, newRefreshToken); err != nil {
			return err
		}

		return nil
	})
}

// DeleteRefreshToken deletes refresh token from the repository.
func (r *Repository) DeleteRefreshToken(ctx context.Context, refreshTokenID id.RefreshToken) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		return r.deleteRefreshToken(dctx, refreshTokenID)
	})
}

// DeleteRefreshTokensByAuthorID deletes all refresh tokens by user id.
func (r *Repository) DeleteRefreshTokensByAuthorID(ctx context.Context, authorID id.Author) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		return sql.Exec(dctx, query.DeleteRefreshTokensByAuthorID, pgx.NamedArgs{
			"user_id": authorID,
		})
	})
}

func (r *Repository) deleteRefreshToken(dctx sql.DataContext, refreshTokenID id.RefreshToken) error {
	err := sql.ExecOne(dctx, query.DeleteRefreshToken, pgx.NamedArgs{
		"id": refreshTokenID,
	})
	if err != nil {
		if sql.IsNotFound(err) {
			return domsession.ErrRefreshTokenNotFound
		}
		return err
	}
	return nil
}
