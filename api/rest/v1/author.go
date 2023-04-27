package v1

import (
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
	domauthor "newsletter-manager-go/domain/author"
	apierrors "newsletter-manager-go/types/errors"
	utilctx "newsletter-manager-go/util/context"
)

// CreateAuthor creates new author.
func (h *Handler) CreateAuthor(_ http.ResponseWriter, r *http.Request, input model.CreateAuthorInput) (*model.CreateAuthorResp, error) {
	createAuthorInput, err := domauthor.NewCreateAuthorInput(
		input.Name,
		input.Email,
		input.Password,
		input.ReferrerID,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create author input").WithPublicMessage(err.Error())
	}
	author, session, err := h.authorService.Create(r.Context(), createAuthorInput)
	if err != nil {
		return nil, fmt.Errorf("creating author: %w", err)
	}
	createAuthorResp := model.CreateAuthorResp{
		Author:  model.FromAuthor(author),
		Session: model.FromSession(session),
	}
	return &createAuthorResp, nil
}

// ReadLoggedAuthor returns existing author.
func (h *Handler) ReadLoggedAuthor(_ http.ResponseWriter, r *http.Request) (*model.Author, error) {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())
	author, err := h.authorService.Read(r.Context(), authorID)
	if err != nil {
		return nil, fmt.Errorf("reading logged author: %w", err)
	}

	authorResp := model.FromAuthor(author)
	return &authorResp, nil
}

// ChangeAuthorPassword changes author password.
func (h *Handler) ChangeAuthorPassword(_ http.ResponseWriter, r *http.Request, input model.ChangeAuthorPasswordInput) error {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())
	if err := h.authorService.ChangePassword(r.Context(), authorID, input.OldPassword, input.NewPassword); err != nil {
		return fmt.Errorf("changing author password: %w", err)
	}
	return nil
}

// ListAuthors returns all existing authors.
func (h *Handler) ListAuthors(_ http.ResponseWriter, r *http.Request) ([]model.Author, error) {
	authors, err := h.authorService.List(r.Context())
	if err != nil {
		return nil, fmt.Errorf("listing authors: %w", err)
	}
	return model.FromAuthors(authors), nil
}
