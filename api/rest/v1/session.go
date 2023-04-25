package v1

import (
	"fmt"
	"net/http"

	"strv-template-backend-go-api/api/rest/v1/model"
)

// CreateSession creates new stateless session.
func (h *Handler) CreateSession(_ http.ResponseWriter, r *http.Request, input model.CreateSessionInput) (*model.CreateUserResp, error) {
	session, user, err := h.sessionService.Create(r.Context(), input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	createSessionResp := model.CreateUserResp{
		User:    model.FromUser(user),
		Session: model.FromSession(session),
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
		Session: model.FromSession(session),
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
