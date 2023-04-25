package middleware

import (
	"errors"
	"net/http"

	httpx "go.strv.io/net/http"

	"go.uber.org/zap"
)

const (
	codePayloadTooLarge = "ERR_PAYLOAD_TOO_LARGE"

	// DefaultByteCountLimit is 4 MiB.
	DefaultByteCountLimit = 4 * 1024 * 1024
)

// LimitBodySize limits size of incoming body. It wraps r.Body by http.MaxBytesReader.
// When limit is reached, reader will return an error.
func LimitBodySize(logger *zap.Logger, byteCountLimit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > byteCountLimit {
				err := errors.New("payload too large")
				if err = httpx.WriteErrorResponse(
					w,
					http.StatusRequestEntityTooLarge,
					httpx.WithError(err),
					httpx.WithErrorCode(codePayloadTooLarge),
				); err != nil {
					logger.With(
						zap.Int("status_code", http.StatusRequestEntityTooLarge),
					).Error("writing http error response", zap.Error(err))
				}
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, byteCountLimit)
			next.ServeHTTP(w, r)
		})
	}
}
