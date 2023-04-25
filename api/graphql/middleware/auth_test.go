package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domsession "strv-template-backend-go-api/domain/session"
	"strv-template-backend-go-api/domain/user"
	"strv-template-backend-go-api/types/id"
	utilctx "strv-template-backend-go-api/util/context"

	httpx "go.strv.io/net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	errTest = errors.New("test error")
)

type mockTokenParser struct {
	mock.Mock
}

func (m *mockTokenParser) ParseAccessToken(data string) (*domsession.AccessToken, error) {
	args := m.Called(data)
	return args.Get(0).(*domsession.AccessToken), args.Error(1)
}

func Test_Authenticate(t *testing.T) {
	userID := id.NewUser()
	userRole := user.RoleUser
	token := "test_access_token"

	tests := []struct {
		name               string
		tokenParser        TokenParser
		handler            http.Handler
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			tokenParser: func() TokenParser {
				accessToken := &domsession.AccessToken{
					Claims: domsession.Claims{
						UserID: userID,
						Custom: domsession.CustomClaims{UserRole: userRole},
					},
				}
				tokenParser := &mockTokenParser{}
				tokenParser.On("ParseAccessToken", token).Return(accessToken, nil)
				return tokenParser
			}(),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Helper()
				uid, ok := utilctx.UserIDFromCtx(r.Context())
				assert.True(t, ok)
				assert.Equal(t, userID, uid)
				urole, ok := utilctx.UserRoleFromCtx(r.Context())
				assert.True(t, ok)
				assert.Equal(t, userRole, urole)
				w.WriteHeader(http.StatusNoContent)
			}),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeader, bearerSchema+token)
				return r
			}(),
			expectedStatusCode: http.StatusNoContent,
			expectedBody:       http.NoBody,
		},
		{
			name: "failure:parse-token",
			tokenParser: func() TokenParser {
				tokenParser := &mockTokenParser{}
				tokenParser.On("ParseAccessToken", token).Return((*domsession.AccessToken)(nil), errTest)
				return tokenParser
			}(),
			handler: nil,
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
				require.NoError(t, err)
				r.Header.Set(authHeader, bearerSchema+token)
				return r
			}(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: errorCodeUnauthorized,
				ErrData: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMiddleware := Authenticate(zap.NewNop(), tt.tokenParser)
			w := httptest.NewRecorder()
			authMiddleware(tt.handler).ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
		})
	}
}

func Test_parseBearerToken(t *testing.T) {
	expected := "test_bearer_token"
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", expected))

	assert.Equal(t, expected, parseBearerToken(header))
	assert.Equal(t, "", parseBearerToken(nil))
}

func assertResponseBody(t *testing.T, expectedBody any, body *bytes.Buffer) {
	t.Helper()

	if expectedBody == http.NoBody {
		assert.Empty(t, body)
		return
	}

	bodyData := body.Bytes()
	expectedBodyData, err := json.Marshal(expectedBody)
	assert.NoError(t, err)
	assert.Equal(t, bytes.TrimSpace(expectedBodyData), bytes.TrimSpace(bodyData))
}
