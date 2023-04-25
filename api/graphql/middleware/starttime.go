package middleware

import (
	"context"
	"net/http"
	"time"

	"strv-template-backend-go-api/util/timesource"
)

type ctxKeyStartTime struct{}

var (
	contextKey = struct {
		startTime ctxKeyStartTime
	}{}
)

// RequestStartTime passes to context current time.
func RequestStartTime(timeSource timesource.TimeSource) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithStartTime(r.Context(), timeSource.Now())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithStartTime passes time to the context.
func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, contextKey.startTime, t)
}

// StartTimeFromCtx gets time from the context.
func StartTimeFromCtx(ctx context.Context) (time.Time, bool) {
	startTime, ok := ctx.Value(contextKey.startTime).(time.Time)
	return startTime, ok
}
