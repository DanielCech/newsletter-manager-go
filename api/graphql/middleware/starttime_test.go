package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"strv-template-backend-go-api/util/timesource"
)

func Test_RequestStartTime(t *testing.T) {
	tests := []struct {
		name    string
		handler http.Handler
		request *http.Request
	}{
		{
			name: "success",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				startTime, ok := StartTimeFromCtx(r.Context())
				assert.True(t, ok)
				assert.Condition(t, func() bool {
					return time.Duration(time.Since(startTime).Milliseconds()) < time.Millisecond*50
				})
			}),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/test", http.NoBody)
				require.NoError(t, err)
				return r
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestStartTimeMiddleware := RequestStartTime(timesource.DefaultTimeSource{})
			requestStartTimeMiddleware(tt.handler).ServeHTTP(httptest.NewRecorder(), tt.request)
		})
	}
}

func Test_WithStartTime(t *testing.T) {
	expected := time.Now()
	ctx := WithStartTime(context.Background(), expected)
	startTime, ok := ctx.Value(contextKey.startTime).(time.Time)
	assert.True(t, ok)
	assert.Equal(t, expected, startTime)
}

func Test_StartTimeFromCtx(t *testing.T) {
	expected := time.Now()
	ctx := context.WithValue(context.Background(), contextKey.startTime, expected)
	startTime, ok := StartTimeFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, expected, startTime)
}
