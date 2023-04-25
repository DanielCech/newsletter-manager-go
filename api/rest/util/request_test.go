package util

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"strv-template-backend-go-api/types"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/types/id"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addURLParam(r *http.Request, key, value string) *http.Request {
	rctx, ok := r.Context().Value(chi.RouteCtxKey).(*chi.Context)
	if !ok {
		rctx = chi.NewRouteContext()
	}
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

type testPasswordObject struct {
	Password types.Password `json:"password"`
}

type testEmailObject struct {
	Email types.Email `json:"email"`
}

type testValidateObject struct {
	FavoriteAnimal string `json:"favAnimal" validate:"required"` // because everybody has to have some favorite animal
}

func Test_ParseRequestBody(t *testing.T) {
	type args struct {
		request *http.Request
		target  any
	}
	tests := []struct {
		name                     string
		args                     args
		expectedErrCode          apierrors.Code
		expectedErrData          any
		expectedErrPublicMessage string
	}{
		{
			name: "success",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					r.Body = io.NopCloser(strings.NewReader(`{"password":"testSecret1"}`))
					return r
				}(),
				target: &testPasswordObject{},
			},
		},
		{
			name: "failure:struct/password",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					r.Body = io.NopCloser(strings.NewReader(`{"password":"testSecret"}`))
					return r
				}(),
				target: &testPasswordObject{},
			},
			expectedErrCode:          apierrors.CodeBadRequest,
			expectedErrPublicMessage: "invalid json body",
		},
		{
			name: "failure:struct/validate-tag-required",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					r.Body = io.NopCloser(strings.NewReader(`{}`))
					return r
				}(),
				target: &testValidateObject{},
			},
			expectedErrCode: apierrors.CodeInvalidBody,
			expectedErrData: map[string]any{
				"invalidFields": []map[string]any{{
					"name": "favAnimal",
				}}},
		},
		{
			name: "failure:struct/email",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					r.Body = io.NopCloser(strings.NewReader(`{"email":""}`))
					return r
				}(),
				target: &testEmailObject{},
			},
			expectedErrCode:          apierrors.CodeBadRequest,
			expectedErrPublicMessage: "invalid json body",
		},
		{
			name: "failure:decoder",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					r.Body = io.NopCloser(strings.NewReader(`{"`))
					return r
				}(),
				target: &testPasswordObject{},
			},
			expectedErrCode:          apierrors.CodeBadRequest,
			expectedErrPublicMessage: "invalid json body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseRequestBody(tt.args.request, tt.args.target)
			if tt.expectedErrCode == "" {
				assert.NoError(t, err)
				return
			}
			e := &apierrors.Error{}
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, tt.expectedErrPublicMessage, e.PublicMessage)
			assert.Equal(t, tt.expectedErrCode, e.Code)
			assert.Equal(t, tt.expectedErrData, e.Data)
		})
	}
}

func Test_GetPathID(t *testing.T) {
	userID := id.NewUser()

	type args struct {
		request   *http.Request
		paramName string
	}
	tests := []struct {
		name           string
		args           args
		expectedUserID id.User
		expectedErr    *apierrors.Error
	}{
		{
			name: "success",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/users/"+userID.String(), http.NoBody)
					require.NoError(t, err)
					r = addURLParam(r, "userId", userID.String())
					return r
				}(),
				paramName: "userId",
			},
			expectedUserID: userID,
			expectedErr:    nil,
		},
		{
			name: "failure:unmarshal-text",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/users/123", http.NoBody)
					require.NoError(t, err)
					return r
				}(),
				paramName: "userId",
			},
			expectedUserID: id.User{},
			expectedErr: apierrors.NewBadRequestError(
				errors.New(`unmarshalling text: parsing "User" id value: invalid UUID length: 0`),
				"",
			).WithPublicMessage(`invalid path id parameter "userId"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetPathID[id.User](tt.args.request, tt.args.paramName)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, result)
				return
			}
			e := &apierrors.Error{}
			assert.ErrorAs(t, err, &e, "unknown error type")
			assert.Equal(t, tt.expectedErr.Error(), e.Error())
			assert.Equal(t, tt.expectedErr.PublicMessage, e.PublicMessage)
			assert.Equal(t, tt.expectedErr.Code, e.Code)
			assert.Equal(t, tt.expectedErr.Data, e.Data)
		})
	}
}

func Test_GetQueryID(t *testing.T) {
	userID := id.NewUser()

	type args struct {
		request   *http.Request
		paramName string
	}
	tests := []struct {
		name           string
		args           args
		expectedUserID *id.User
		expectedErr    *apierrors.Error
	}{
		{
			name: "success",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test?userId="+userID.String(), http.NoBody)
					require.NoError(t, err)
					return r
				}(),
				paramName: "userId",
			},
			expectedUserID: &userID,
			expectedErr:    nil,
		},
		{
			name: "success:no-query-parameter-set",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test", http.NoBody)
					require.NoError(t, err)
					return r
				}(),
				paramName: "userId",
			},
			expectedUserID: &id.User{},
			expectedErr:    nil,
		},
		{
			name: "failure:unmarshal-text",
			args: args{
				request: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "/test?userId=abc-123", http.NoBody)
					require.NoError(t, err)
					return r
				}(),
				paramName: "userId",
			},
			expectedUserID: nil,
			expectedErr: apierrors.NewBadRequestError(
				errors.New(`unmarshalling text: parsing "User" id value: invalid UUID length: 7`),
				"",
			).WithPublicMessage(`invalid query id parameter "userId"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetQueryID[id.User](tt.args.request, tt.args.paramName)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, result)
				return
			}
			e := &apierrors.Error{}
			assert.ErrorAs(t, err, &e, "unknown error type")
			assert.Equal(t, tt.expectedErr.Error(), e.Error())
			assert.Equal(t, tt.expectedErr.PublicMessage, e.PublicMessage)
			assert.Equal(t, tt.expectedErr.Code, e.Code)
			assert.Equal(t, tt.expectedErr.Data, e.Data)
		})
	}
}
