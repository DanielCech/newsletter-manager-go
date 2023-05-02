package author

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"newsletter-manager-go/database/sql"
	domauthor "newsletter-manager-go/domain/author"
	"newsletter-manager-go/domain/author/postgres/query"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Repository represents author data layer.
// Every model that is returned is converted to domain model using author factory.
type Repository struct {
	dataSource    sql.DataSource
	authorFactory domauthor.Factory
}

// NewRepository returns new instance of a author repository.
func NewRepository(dataSource sql.DataSource, authorFactory domauthor.Factory) (*Repository, error) {
	if dataSource == nil {
		return nil, errors.New("invalid data source")
	}
	return &Repository{
		dataSource:    dataSource,
		authorFactory: authorFactory,
	}, nil
}

// Create creates author in the repository.
func (r *Repository) Create(ctx context.Context, author *domauthor.Author) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		err := sql.Exec(dctx, query.Create, pgx.NamedArgs{
			"id":            author.ID,
			"name":          author.Name,
			"email":         author.Email,
			"password_hash": author.PasswordHash,
			"created_at":    author.CreatedAt,
			"updated_at":    author.UpdatedAt,
		})
		if err != nil {
			//if sql.IsUnique(err, query.ConstraintUniqueAuthorEmail) {
			//	return domauthor.ErrAuthorEmailAlreadyExists
			//}
			//if sql.IsForeignKey(err, query.ConstraintReferrerID) {
			//	return domauthor.ErrReferrerNotFound
			//}
			return err
		}
		return nil
	})
}

// Read reads the author from the repository.
func (r *Repository) Read(ctx context.Context, authorID id.Author) (*domauthor.Author, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domauthor.Author, error) {
		return r.read(dctx, authorID)
	})
}

func (r *Repository) read(dctx sql.DataContext, authorID id.Author) (*domauthor.Author, error) {
	author, err := sql.ReadValue[author](dctx, query.Read, pgx.NamedArgs{
		"id": authorID,
	})
	if err != nil {
		if sql.IsNotFound(err) {
			return nil, domauthor.ErrAuthorNotFound
		}
		return nil, err
	}
	return author.ToAuthor(r.authorFactory), nil
}

// ReadByEmail reads the author by email from the repository.
func (r *Repository) ReadByEmail(ctx context.Context, email types.Email) (*domauthor.Author, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domauthor.Author, error) {
		author, err := sql.ReadValue[author](dctx, query.ReadByEmail, pgx.NamedArgs{
			"email": email,
		})
		if err != nil {
			if sql.IsNotFound(err) {
				return nil, domauthor.ErrAuthorNotFound
			}
			return nil, err
		}
		return author.ToAuthor(r.authorFactory), nil
	})
}

// List lists authors from the repository.
func (r *Repository) List(ctx context.Context) ([]domauthor.Author, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) ([]domauthor.Author, error) {
		dbAuthors, err := sql.List[author](dctx, query.List)
		if err != nil {
			return nil, err
		}
		authors := make([]domauthor.Author, 0, len(dbAuthors))
		for _, u := range dbAuthors {
			authors = append(authors, *u.ToAuthor(r.authorFactory))
		}
		return authors, nil
	})
}

// Update reads the author, calls external update function and updates the author in the repository.
func (r *Repository) Update(ctx context.Context, authorID id.Author, updateFn domauthor.UpdateFunc) error {
	return sql.WithTransaction(ctx, r.dataSource, func(dctx sql.DataContext) error {
		originalAuthor, err := r.read(dctx, authorID)
		if err != nil {
			return err
		}

		newAuthor, err := updateFn(originalAuthor)
		if err != nil {
			return err
		}

		return sql.ExecOne(dctx, query.Update, pgx.NamedArgs{
			"id":            newAuthor.ID,
			"name":          newAuthor.Name,
			"email":         newAuthor.Email,
			"password_hash": newAuthor.PasswordHash,
			"created_at":    newAuthor.CreatedAt,
			"updated_at":    newAuthor.UpdatedAt,
		})
	})
}
