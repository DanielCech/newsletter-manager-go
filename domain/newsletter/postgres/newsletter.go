package newsletter

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"newsletter-manager-go/database/sql"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/domain/newsletter/postgres/query"
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Repository represents newsletter data layer.
// Every model that is returned is converted to domain model using newsletter factory.
type Repository struct {
	dataSource        sql.DataSource
	newsletterFactory domnewsletter.Factory
}

// NewRepository returns new instance of a newsletter repository.
func NewRepository(dataSource sql.DataSource, newsletterFactory domnewsletter.Factory) (*Repository, error) {
	if dataSource == nil {
		return nil, errors.New("invalid data source")
	}
	return &Repository{
		dataSource:        dataSource,
		newsletterFactory: newsletterFactory,
	}, nil
}

// Create creates newsletter in the repository.
func (r *Repository) Create(ctx context.Context, newsletter *domnewsletter.Newsletter) error {
	return sql.WithConnection(ctx, r.dataSource, func(dctx sql.DataContext) error {
		err := sql.Exec(dctx, query.Create, pgx.NamedArgs{
			"id":            newsletter.ID,
			"referrer_id":   newsletter.ReferrerID,
			"name":          newsletter.Name,
			"email":         newsletter.Email,
			"password_hash": newsletter.PasswordHash,
			"role":          newsletter.Role,
			"created_at":    newsletter.CreatedAt,
			"updated_at":    newsletter.UpdatedAt,
		})
		if err != nil {
			//if sql.IsUnique(err, query.ConstraintUniqueNewsletterEmail) {
			//	return domnewsletter.ErrNewsletterEmailAlreadyExists
			//}
			//if sql.IsForeignKey(err, query.ConstraintReferrerID) {
			//	return domnewsletter.ErrReferrerNotFound
			//}
			return err
		}
		return nil
	})
}

// Read reads the newsletter from the repository.
func (r *Repository) Read(ctx context.Context, newsletterID id.Newsletter) (*domnewsletter.Newsletter, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domnewsletter.Newsletter, error) {
		return r.read(dctx, newsletterID)
	})
}

func (r *Repository) read(dctx sql.DataContext, newsletterID id.Newsletter) (*domnewsletter.Newsletter, error) {
	newsletter, err := sql.ReadValue[newsletter](dctx, query.Read, pgx.NamedArgs{
		"id": newsletterID,
	})
	if err != nil {
		if sql.IsNotFound(err) {
			return nil, domnewsletter.ErrNewsletterNotFound
		}
		return nil, err
	}
	return newsletter.ToNewsletter(r.newsletterFactory), nil
}

// ReadByEmail reads the newsletter by email from the repository.
func (r *Repository) ReadByEmail(ctx context.Context, email types.Email) (*domnewsletter.Newsletter, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) (*domnewsletter.Newsletter, error) {
		newsletter, err := sql.ReadValue[newsletter](dctx, query.ReadByEmail, pgx.NamedArgs{
			"email": email,
		})
		if err != nil {
			if sql.IsNotFound(err) {
				return nil, domnewsletter.ErrNewsletterNotFound
			}
			return nil, err
		}
		return newsletter.ToNewsletter(r.newsletterFactory), nil
	})
}

// List lists newsletters from the repository.
func (r *Repository) List(ctx context.Context) ([]domnewsletter.Newsletter, error) {
	return sql.WithConnectionResult(ctx, r.dataSource, func(dctx sql.DataContext) ([]domnewsletter.Newsletter, error) {
		dbNewsletters, err := sql.List[newsletter](dctx, query.List)
		if err != nil {
			return nil, err
		}
		newsletters := make([]domnewsletter.Newsletter, 0, len(dbNewsletters))
		for _, u := range dbNewsletters {
			newsletters = append(newsletters, *u.ToNewsletter(r.newsletterFactory))
		}
		return newsletters, nil
	})
}

// Update reads the newsletter, calls external update function and updates the newsletter in the repository.
func (r *Repository) Update(ctx context.Context, newsletterID id.Newsletter, updateFn domnewsletter.UpdateFunc) error {
	return sql.WithTransaction(ctx, r.dataSource, func(dctx sql.DataContext) error {
		originalNewsletter, err := r.read(dctx, newsletterID)
		if err != nil {
			return err
		}

		newNewsletter, err := updateFn(originalNewsletter)
		if err != nil {
			return err
		}

		return sql.ExecOne(dctx, query.Update, pgx.NamedArgs{
			"id":            newNewsletter.ID,
			"referrer_id":   newNewsletter.ReferrerID,
			"name":          newNewsletter.Name,
			"email":         newNewsletter.Email,
			"password_hash": newNewsletter.PasswordHash,
			"role":          newNewsletter.Role,
			"created_at":    newNewsletter.CreatedAt,
			"updated_at":    newNewsletter.UpdatedAt,
		})
	})
}
