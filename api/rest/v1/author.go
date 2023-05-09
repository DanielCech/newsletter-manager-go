package v1

import (
	"fmt"
	"net/http"
	domauthor "newsletter-manager-go/domain/author"
	apierrors "newsletter-manager-go/types/errors"

	"newsletter-manager-go/api/rest/v1/model"
	utilctx "newsletter-manager-go/util/context"
)

func (h *Handler) AuthorSignUp(_ http.ResponseWriter, r *http.Request, input model.AuthorSignUpInput) (*model.CreateAuthorResp, error) {
	createAuthorInput, err := domauthor.NewCreateAuthorInput(
		input.Name,
		input.Email,
		input.Password,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create author input").WithPublicMessage(err.Error())
	}
	domauthor, session, err := h.authorService.Create(r.Context(), createAuthorInput)
	if err != nil || domauthor == nil {
		return nil, fmt.Errorf("creating author: %w", err)
	}

	author := model.FromDomainAuthor(domauthor)

	createAuthorResp := model.CreateAuthorResp{
		Author:  &author,
		Session: model.FromDomainSession(session),
	}
	return &createAuthorResp, nil
}

func (h *Handler) AuthorSignIn(_ http.ResponseWriter, r *http.Request, input model.AuthorSignInInput) (*model.FullAuthor, error) {
	// TODO
	return nil, nil
}

func (h *Handler) RefreshToken(_ http.ResponseWriter, r *http.Request, input model.RefreshTokenInput) (*model.FullAuthor, error) {
	// TODO
	return nil, nil
}

func (h *Handler) ListAuthors(_ http.ResponseWriter, r *http.Request) ([]model.Author, error) {
	// TODO
	// return model.FromAuthors(authors), nil
	return nil, nil
}

func (h *Handler) GetAuthor(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) (*model.Author, error) {
	// TODO
	return nil, nil
}

func (h *Handler) UpdateAuthor(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) (*model.Author, error) {
	// TODO
	return nil, nil
}

func (h *Handler) DeleteAuthor(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) error {
	// TODO
	return nil
}

// ReadLoggedAuthor returns existing author.
func (h *Handler) ReadLoggedAuthor(_ http.ResponseWriter, r *http.Request) (*model.Author, error) {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())
	author, err := h.authorService.Read(r.Context(), authorID)
	if err != nil {
		return nil, fmt.Errorf("reading logged author: %w", err)
	}

	authorResp := model.FromDomainAuthor(author)
	return &authorResp, nil
}
