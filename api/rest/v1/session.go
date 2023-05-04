package v1

import (
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/v1/model"
)

// CreateSession creates new stateless session.
func (h *Handler) CreateSession(_ http.ResponseWriter, r *http.Request, input model.CreateSessionInput) (*model.CreateAuthorResponse, error) {
	dsession, dauthor, err := h.sessionService.Create(r.Context(), input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	author := model.FromDomainAuthor(dauthor)
	session := model.FromDomainSession(dsession)

	createSessionResp := model.CreateAuthorResponse{
		Author:  &author,
		Session: session,
	}
	return &createSessionResp, nil
}

// RefreshSession refreshes existing stateless session.
func (h *Handler) RefreshSession(_ http.ResponseWriter, r *http.Request, input model.RefreshSessionInput) (*model.RefreshSessionResp, error) {
	session, err := h.sessionService.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refreshing session: %w", err)
	}
	refreshSessionResp := model.RefreshSessionResp{
		Session: model.FromDomainSession(session),
	}
	return &refreshSessionResp, nil
}

// DestroySession deletes existing stateless session.
func (h *Handler) DestroySession(_ http.ResponseWriter, r *http.Request, input model.DestroySessionInput) error {
	if err := h.sessionService.Destroy(r.Context(), input.RefreshToken); err != nil {
		return fmt.Errorf("destroying session: %w", err)
	}
	return nil
}
