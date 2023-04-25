package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "strv-template-backend-go-api/types/errors"

	httpx "go.strv.io/net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_LimitBodySize(t *testing.T) {
	type args struct {
		byteCountLimit int64
	}
	tests := []struct {
		name               string
		args               args
		handler            http.Handler
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			args: args{byteCountLimit: DefaultByteCountLimit},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}),
			request: func() *http.Request {
				body := io.NopCloser(strings.NewReader(`{"field1":"value"}`))
				r, err := http.NewRequest(http.MethodPost, "/test", body)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusNoContent,
			expectedBody:       http.NoBody,
		},
		{
			name: "failure:payload-too-large",
			args: args{byteCountLimit: 5},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}),
			request: func() *http.Request {
				d := `{"field1":"value"}`
				body := io.NopCloser(strings.NewReader(d))
				r, err := http.NewRequest(http.MethodPost, "/test", body)
				require.NoError(t, err)
				r.ContentLength = int64(len(d))
				return r
			}(),
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			expectedBody: httpx.ErrorResponseOptions{
				ErrCode: string(apierrors.CodePayloadTooLarge),
				ErrData: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limitBodySizeMiddleware := LimitBodySize(zap.NewNop(), tt.args.byteCountLimit)
			w := httptest.NewRecorder()
			limitBodySizeMiddleware(tt.handler).ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
		})
	}
}
