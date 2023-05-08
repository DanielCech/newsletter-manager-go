package middleware

import (
	"errors"
	"net/http"
	httputil "newsletter-manager-go/api/rest/util"
	domsession "newsletter-manager-go/domain/session"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/util"
	utilctx "newsletter-manager-go/util/context"
	"strings"

	"go.uber.org/zap"
)

const (
	authHeader   = "Authorization"
	bearerSchema = "Bearer "
)

var (
	ErrMissingToken         = errors.New("missing token")
	ErrMissingUserRole      = errors.New("missing user role")
	ErrInsufficientUserRole = errors.New("insufficient user role")
)

// TokenParser is an interface for parsing incoming bearer tokens.
type TokenParser interface {
	ParseAccessToken(data string) (*domsession.AccessToken, error)
}

// Authenticate parses bearer token from authorization header.
// Custom claims parsed from access token are passed to context.
func Authenticate(logger *zap.Logger, tokenParser TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := parseBearerToken(r.Header)
			if len(token) == 0 {
				httputil.WriteErrorResponse(
					r.Context(),
					util.WithCtx(r.Context(), logger),
					w,
					apierrors.NewUnauthorizedError(ErrMissingToken, "parsing auth token from http header"),
				)
				return
			}

			accessToken, err := tokenParser.ParseAccessToken(token)
			if err != nil {
				httputil.WriteErrorResponse(
					r.Context(),
					util.WithCtx(r.Context(), logger),
					w,
					apierrors.NewUnauthorizedError(err, "parsing access token"),
				)
				return
			}

			ctx := utilctx.WithAuthorID(r.Context(), accessToken.Claims.AuthorID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseBearerToken(h http.Header) string {
	if h == nil {
		return ""
	}
	return strings.TrimPrefix(h.Get(authHeader), bearerSchema)
}
