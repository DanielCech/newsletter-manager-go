package v1

import (
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	apierrors "newsletter-manager-go/types/errors"
)

// CreateNewsletter creates new user.
func (h *Handler) CreateNewsletter(_ http.ResponseWriter, r *http.Request, input model.CreateNewsletterInput) (*model.CreateNewsletterResp, error) {
	createNewsletterInput, err := domnewsletter.NewCreateNewsletterInput(
		input.Name,
		input.Email,
		input.Password,
		input.ReferrerID,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create user input").WithPublicMessage(err.Error())
	}
	newsletter, session, err := h.userService.Create(r.Context(), createNewsletterInput)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	createNewsletterResp := model.CreateNewsletterResp{
		Newsletter: model.FromNewsletter(newsletter),
		Session:    model.FromSession(session),
	}
	return &createNewsletterResp, nil
}

// ListNewsletters returns all existing users.
func (h *Handler) ListNewsletters(_ http.ResponseWriter, r *http.Request) ([]model.Newsletter, error) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return model.FromNewsletters(users), nil
}
