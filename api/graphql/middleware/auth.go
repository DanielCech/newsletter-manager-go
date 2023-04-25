package middleware

import (
	"net/http"
	"strings"

	domsession "strv-template-backend-go-api/domain/session"
	utilctx "strv-template-backend-go-api/util/context"

	httpx "go.strv.io/net/http"

	"go.uber.org/zap"
)

const (
	authHeader   = "Authorization"
	bearerSchema = "Bearer "

	errorCodeUnauthorized = "ERROR_UNAUTHORIZED"
)

// TokenParser is an interface for parsing incoming bearer tokens.
type TokenParser interface {
	ParseAccessToken(data string) (*domsession.AccessToken, error)
}

// Authenticate parses bearer token from authorization header.
// Custom claims parsed from access token are passed to context.
// Bearer token is optional. Subsequent authorization should be done by graphql directive.
func Authenticate(logger *zap.Logger, tokenParser TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := parseBearerToken(r.Header)
			if len(token) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			accessToken, err := tokenParser.ParseAccessToken(token)
			if err != nil {
				if err = httpx.WriteErrorResponse(
					w,
					http.StatusUnauthorized,
					httpx.WithError(err),
					httpx.WithErrorCode(errorCodeUnauthorized),
				); err != nil {
					logger.With(
						zap.Int("status_code", http.StatusUnauthorized),
					).Error("writing http error response", zap.Error(err))
				}
				return
			}

			ctx := utilctx.WithUserID(r.Context(), accessToken.Claims.UserID)
			ctx = utilctx.WithUserRole(ctx, accessToken.Claims.Custom.UserRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseBearerToken(h http.Header) string {
	return strings.TrimPrefix(h.Get(authHeader), bearerSchema)
}
