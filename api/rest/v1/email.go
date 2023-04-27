package v1

import (
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	apierrors "newsletter-manager-go/types/errors"
)

// CreateEmail creates new email.
func (h *Handler) CreateEmail(_ http.ResponseWriter, r *http.Request, input model.CreateEmailInput) (*model.CreateEmailResp, error) {
	createEmailInput, err := domnewsletter.NewCreateEmailInput(
		input.Name,
		input.Email,
		input.Password,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create email input").WithPublicMessage(err.Error())
	}
	email, session, err := h.newsletterService.Create(r.Context(), createEmailInput)
	if err != nil {
		return nil, fmt.Errorf("creating email: %w", err)
	}
	createEmailResp := model.CreateEmailResp{
		Email:   model.FromEmail(email),
		Session: model.FromSession(session),
	}
	return &createEmailResp, nil
}

// ListEmails returns all existing emails.
func (h *Handler) ListEmails(_ http.ResponseWriter, r *http.Request) ([]model.Email, error) {
	emails, err := h.newsletterService.List(r.Context())
	if err != nil {
		return nil, fmt.Errorf("listing emails: %w", err)
	}
	return model.FromEmails(emails), nil
}
