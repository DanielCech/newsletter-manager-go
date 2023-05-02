package v1

import (
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
)

func (h *Handler) AuthorSignUp(_ http.ResponseWriter, r *http.Request, input model.AuthorSignUpInput) (*model.FullAuthor, error) {
	// TODO
	return nil, nil
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
	//return model.FromAuthors(authors), nil
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
