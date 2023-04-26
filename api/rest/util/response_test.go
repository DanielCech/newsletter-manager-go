package util

import (
	"errors"
	"net/http"
	"testing"

	apierrors "newsletter-manager-go/types/errors"

	httpx "go.strv.io/net/http"

	"github.com/stretchr/testify/assert"
)

var (
	errTest    = errors.New("test error")
	errMessage = "test error message"
)

func Test_ToStatusCodeWithOptions(t *testing.T) {
	errData := map[string]any{
		"testKey": "testValue",
	}

	type args struct {
		err error
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
		expectedOpts       []httpx.ErrorResponseOption
	}{
		{
			name: "default-unknown-error-type",
			args: args{
				err: errTest,
			},
			expectedStatusCode: defaultErrorStatusCode,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeUnknown)),
			},
		},
		{
			name: "default-api-error-status-code",
			args: func() args {
				err := apierrors.NewError(errTest, "", "")
				err = err.WithData(errData)
				return args{err: err}
			}(),
			expectedStatusCode: defaultErrorStatusCode,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(""),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(errData),
			},
		},
		{
			name: "bad-request",
			args: args{
				err: apierrors.NewBadRequestError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeBadRequest)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "unauthorized",
			args: args{
				err: apierrors.NewUnauthorizedError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeUnauthorized)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "forbidden",
			args: args{
				err: apierrors.NewForbiddenError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusForbidden,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeForbidden)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "not-found",
			args: args{
				err: apierrors.NewNotFoundError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusNotFound,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeNotFound)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "already-exists",
			args: args{
				err: apierrors.NewAlreadyExistsError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusConflict,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeAlreadyExists)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "expired",
			args: args{
				err: apierrors.NewExpiredError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusGone,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeExpired)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "payload-too-large",
			args: args{
				err: apierrors.NewPayloadTooLargeError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodePayloadTooLarge)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "invalid-body",
			args: args{
				err: apierrors.NewInvalidBodyError(errTest, errMessage),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeInvalidBody)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
		{
			name: "unknown",
			args: args{
				err: apierrors.NewUnknownError(errTest, errMessage),
			},
			expectedStatusCode: defaultErrorStatusCode,
			expectedOpts: []httpx.ErrorResponseOption{
				httpx.WithError(errTest),
				httpx.WithErrorCode(string(apierrors.CodeUnknown)),
				httpx.WithErrorMessage(""),
				httpx.WithErrorData(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, opts := ToStatusCodeWithOptions(tt.args.err)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Len(t, opts, len(tt.expectedOpts))
		})
	}
}
