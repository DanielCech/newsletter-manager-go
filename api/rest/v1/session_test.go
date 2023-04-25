package v1

import (
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
	"strv-template-backend-go-api/util/timesource"

	httpx "go.strv.io/net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newSession(t *testing.T, customClaims domsession.Claims) *domsession.Session {
	t.Helper()
	factory, err := domsession.NewFactory([]byte("abc123"), timesource.DefaultTimeSource{}, time.Hour, time.Hour)
	require.NoError(t, err)
	session, err := factory.NewSession(customClaims)
	require.NoError(t, err)
	return session
}

func Test_Handler_CreateSession(t *testing.T) {
	password := types.Password("Topsecret1")
	user := newUser()
	session := newSession(t, domsession.Claims{
		UserID: user.ID,
		Custom: domsession.CustomClaims{UserRole: user.Role},
	})
	createSessionResp := model.CreateSessionResp{
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
				mocks.sessionService.On("Create", user.Email, password).Return(session, user, nil)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"email":"%s",
					"password":"%s"
				}`,
					user.Email,
					string(password),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/native", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusCreated,
			expectedBody:       createSessionResp,
		},
		{
			name: "failure:create",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Create", user.Email, password).Return((*domsession.Session)(nil), (*domuser.User)(nil), serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{
					"email":"%s",
					"password":"%s"
				}`,
					user.Email,
					string(password),
				)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/native", body)
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
				r, err := http.NewRequest(http.MethodPost, "/sessions/native", http.NoBody)
				require.NoError(t, err)
				r.Body = io.NopCloser(strings.NewReader(`{"`))
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

func Test_Handler_RefreshSession(t *testing.T) {
	inputRefreshToken := id.RefreshToken("5asd4a6d4a36d45as36da")
	session := newSession(t, domsession.Claims{
		UserID: id.NewUser(),
		Custom: domsession.CustomClaims{UserRole: domuser.RoleUser},
	})
	refreshSessionResp := model.RefreshSessionResp{
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
				mocks.sessionService.On("Refresh", inputRefreshToken).Return(session, nil)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{"refreshToken":"%s"}`, inputRefreshToken)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/refresh", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusCreated,
			expectedBody:       refreshSessionResp,
		},
		{
			name: "failure:refresh",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Refresh", inputRefreshToken).Return((*domsession.Session)(nil), serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{"refreshToken":"%s"}`, inputRefreshToken)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/refresh", body)
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
				r, err := http.NewRequest(http.MethodPost, "/sessions/refresh", body)
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

func Test_Handler_DestroySession(t *testing.T) {
	inputRefreshToken := id.RefreshToken("5asd4a6d4a36d45as36da")
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
				mocks.sessionService.On("Destroy", inputRefreshToken).Return(nil)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{"refreshToken":"%s"}`, inputRefreshToken)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/destroy", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusNoContent,
			expectedBody:       http.NoBody,
		},
		{
			name: "failure:destroy",
			mocks: func() mocks {
				mocks := newMocks()
				mocks.sessionService.On("Destroy", inputRefreshToken).Return(serviceErr)
				return mocks
			}(),
			request: func() *http.Request {
				data := fmt.Sprintf(`{"refreshToken":"%s"}`, inputRefreshToken)
				body := io.NopCloser(strings.NewReader(data))
				r, err := http.NewRequest(http.MethodPost, "/sessions/destroy", body)
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
				r, err := http.NewRequest(http.MethodPost, "/sessions/destroy", body)
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
