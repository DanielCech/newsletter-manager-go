package middleware

import (
	"errors"
	"net/http"

	httputil "newsletter-manager-go/api/rest/util"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/util"

	"go.uber.org/zap"
)

const (
	// DefaultByteCountLimit is 4 MiB.
	DefaultByteCountLimit = 4 * 1024 * 1024
)

// LimitBodySize limits size of incoming body. It wraps r.Body by http.MaxBytesReader.
// When limit is reached, reader will return an error.
func LimitBodySize(logger *zap.Logger, byteCountLimit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > byteCountLimit {
				httputil.WriteErrorResponse(
					r.Context(),
					util.WithCtx(r.Context(), logger),
					w,
					apierrors.NewPayloadTooLargeError(errors.New("payload too large"), ""),
				)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, byteCountLimit)
			next.ServeHTTP(w, r)
		})
	}
}
