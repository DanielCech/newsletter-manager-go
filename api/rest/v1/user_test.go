package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"strv-template-backend-go-api/api/rest/v1/model"
	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"
	utilctx "strv-template-backend-go-api/util/context"

	httpx "go.strv.io/net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	authHeaderKey          = "Authorization"
	testAccessToken        = "xyzAccessToken"
	testAuthorizationValue = "Bearer " + testAccessToken
)

func newUser() *domuser.User {
	referrerID := id.NewUser()
	now := time.Now()
	return &domuser.User{
		ID:           id.NewUser(),
		ReferrerID:   &referrerID,
		Name:         "Jozko Dlouhy",
		Email:        "jozko.dlouhy@gmail.com",
		PasswordHash: []byte("45asdad4as25as"),
		Role:         domuser.RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func Test_Handler_CreateUser(t *testing.T) {
	user := newUser()
	session := newSession(t, domsession.Claims{
		UserID: user.ID,
		Custom: domsession.CustomClaims{UserRole: user.Role},
	})
	createUserInput := domuser.CreateUserInput{
		Name:       user.Name,
		Email:      user.Email,
		Password:   types.Password("Topsecret1"),
		ReferrerID: user.ReferrerID,
	}
	createUserResp := model.CreateUserResp{
		User:    model.FromUser(user),
		Session: model.FromSession(session),
	}
	serviceErr := apierrors.NewUnknownError(errTest, "")

	tests := []struct {
		name               string
		mocks              mocks
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Create", createUserInput).Return(user, session, nil)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"email":"%s",
					"password":"%s",
					"name":"%s",
					"referrerId":"%s"
				}`, createUserInput.Email,
					string(createUserInput.Password),
					createUserInput.Name,
					createUserInput.ReferrerID.String(),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/users/register", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusCreated,
			expectedBody:       createUserResp,
		},
		{
			name: "failure:create-user",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Create", createUserInput).Return((*domuser.User)(nil), (*domsession.Session)(nil), serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"email":"%s",
					"password":"%s",
					"name":"%s",
					"referrerId":"%s"
				}`, createUserInput.Email,
					string(createUserInput.Password),
					createUserInput.Name,
					createUserInput.ReferrerID.String(),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/users/register", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: string(apierrors.CodeUnknown),
			},
		},
		{
			name:  "failure:parse-request-body",
			mocks: newMocks(),
			request: func() *http.Request {
				body := io.NopCloser(strings.NewReader(`{"`))
				r, err := http.NewRequest(http.MethodPost, "/users/register", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode:    string(apierrors.CodeBadRequest),
				ErrMessage: "invalid json body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mocks.userService, tt.mocks.sessionService, tt.mocks.tokenParser, zap.NewNop())
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Handler_ReadLoggedUser(t *testing.T) {
	user := newUser()
	userResp := model.FromUser(user)
	serviceErr := apierrors.NewUnknownError(errTest, "")

	tests := []struct {
		name               string
		mocks              mocks
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Read", user.ID).Return(user, nil)
				return mocks
			}(),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/users/me", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				ctx := utilctx.WithUserID(context.Background(), user.ID)
				return r.WithContext(ctx)
			}(),
			expectedStatusCode: http.StatusOK,
			expectedBody:       userResp,
		},
		{
			name: "failure:read",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("Read", user.ID).Return((*domuser.User)(nil), serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/users/me", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				ctx := utilctx.WithUserID(context.Background(), user.ID)
				return r.WithContext(ctx)
			}(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: string(apierrors.CodeUnknown),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mocks.userService, tt.mocks.sessionService, tokenParserAlwaysUser(user.ID), zap.NewNop())
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Handler_ChangeUserPassword(t *testing.T) {
	userID := id.NewUser()
	oldPassword := types.Password("Topsecret1")
	newPassword := types.Password("Topsecret2")
	serviceErr := apierrors.NewUnknownError(errTest, "")

	tests := []struct {
		name               string
		mocks              mocks
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ChangePassword", userID, oldPassword, newPassword).Return(nil)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"oldPassword":"%s",
					"newPassword":"%s"
				}`,
					string(oldPassword),
					string(newPassword),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPatch, "/users/change-password", body)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				ctx := utilctx.WithUserID(context.Background(), userID)
				return r.WithContext(ctx)
			}(),
			expectedStatusCode: http.StatusNoContent,
			expectedBody:       http.NoBody,
		},
		{
			name: "failure:change-password",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("ChangePassword", userID, oldPassword, newPassword).Return(serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"oldPassword":"%s",
					"newPassword":"%s"
				}`,
					string(oldPassword),
					string(newPassword),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPatch, "/users/change-password", body)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				ctx := utilctx.WithUserID(context.Background(), userID)
				return r.WithContext(ctx)
			}(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: string(apierrors.CodeUnknown),
			},
		},
		{
			name:  "failure:parse-request-body",
			mocks: newMocks(),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodPatch, "/users/change-password", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				r.Body = io.NopCloser(strings.NewReader(`{"`))
				ctx := utilctx.WithUserID(context.Background(), userID)
				return r.WithContext(ctx)
			}(),
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode:    string(apierrors.CodeBadRequest),
				ErrMessage: "invalid json body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mocks.userService, tt.mocks.sessionService, tokenParserAlwaysUser(userID), zap.NewNop())
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
			tt.mocks.assertExpectations(t)
		})
	}
}

func Test_Handler_ListUsers(t *testing.T) {
	users := []domuser.User{*newUser(), *newUser()}
	usersResp := model.FromUsers(users)
	serviceErr := apierrors.NewUnknownError(errTest, "")

	tests := []struct {
		name               string
		mocks              mocks
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("List").Return(users, nil)
				return mocks
			}(),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/users", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				return r
			}(),
			expectedStatusCode: http.StatusOK,
			expectedBody:       usersResp,
		},
		{
			name: "failure:list-users",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.userService.On("List").Return([]domuser.User(nil), serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/users", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeaderKey, testAuthorizationValue)
				return r
			}(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: string(apierrors.CodeUnknown),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mocks.userService, tt.mocks.sessionService, tokenParserAlwaysAdmin(), zap.NewNop())
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
			tt.mocks.assertExpectations(t)
		})
	}
}

func tokenParserAlwaysUser(userID id.User) *mockTokenParser {
	mock := &mockTokenParser{}
	mock.On("ParseAccessToken", testAccessToken).Return(&domsession.AccessToken{
		Claims: domsession.Claims{
			UserID: userID,
			Custom: domsession.CustomClaims{UserRole: domuser.RoleUser},
		},
	}, nil).Maybe()
	return mock
}

func tokenParserAlwaysAdmin() *mockTokenParser {
	mock := &mockTokenParser{}
	mock.On("ParseAccessToken", testAccessToken).Return(&domsession.AccessToken{
		Claims: domsession.Claims{
			UserID: id.NewUser(),
			Custom: domsession.CustomClaims{UserRole: domuser.RoleAdmin},
		},
	}, nil).Maybe()
	return mock
}
