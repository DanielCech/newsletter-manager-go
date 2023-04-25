package middleware

import (
	"net/http"
)

// NoCacheHeaders sets Cache-Control and Pragma headers to not cache.
// It is recommended when returning session data to not let clients cache them.
func NoCacheHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Pragma", "no-cache")
			next.ServeHTTP(w, r)
		})
	}
}
