package v1

import (
	"fmt"
	"net/http"

	"strv-template-backend-go-api/api/rest/v1/model"
	domuser "strv-template-backend-go-api/domain/user"
	apierrors "strv-template-backend-go-api/types/errors"
	utilctx "strv-template-backend-go-api/util/context"
)

// CreateUser creates new user.
func (h *Handler) CreateUser(_ http.ResponseWriter, r *http.Request, input model.CreateUserInput) (*model.CreateUserResp, error) {
	createUserInput, err := domuser.NewCreateUserInput(
		input.Name,
		input.Email,
		input.Password,
		input.ReferrerID,
	)
	if err != nil {
		return nil, apierrors.NewInvalidBodyError(err, "new create user input").WithPublicMessage(err.Error())
	}
	user, session, err := h.userService.Create(r.Context(), createUserInput)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	createUserResp := model.CreateUserResp{
		User:    model.FromUser(user),
		Session: model.FromSession(session),
	}
	return &createUserResp, nil
}

// ReadLoggedUser returns existing user.
func (h *Handler) ReadLoggedUser(_ http.ResponseWriter, r *http.Request) (*model.User, error) {
	userID, _ := utilctx.UserIDFromCtx(r.Context())
	user, err := h.userService.Read(r.Context(), userID)
	if err != nil {
		return nil, fmt.Errorf("reading logged user: %w", err)
	}

	userResp := model.FromUser(user)
	return &userResp, nil
}

// ChangeUserPassword changes user password.
func (h *Handler) ChangeUserPassword(_ http.ResponseWriter, r *http.Request, input model.ChangeUserPasswordInput) error {
	userID, _ := utilctx.UserIDFromCtx(r.Context())
	if err := h.userService.ChangePassword(r.Context(), userID, input.OldPassword, input.NewPassword); err != nil {
		return fmt.Errorf("changing user password: %w", err)
	}
	return nil
}

// ListUsers returns all existing users.
func (h *Handler) ListUsers(_ http.ResponseWriter, r *http.Request) ([]model.User, error) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return model.FromUsers(users), nil
}
