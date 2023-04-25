package middleware

import (
	"errors"
	"net/http"
	"strings"

	httputil "strv-template-backend-go-api/api/rest/util"
	domsession "strv-template-backend-go-api/domain/session"
	domuser "strv-template-backend-go-api/domain/user"
	apierrors "strv-template-backend-go-api/types/errors"
	"strv-template-backend-go-api/util"
	utilctx "strv-template-backend-go-api/util/context"

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

			ctx := utilctx.WithUserID(r.Context(), accessToken.Claims.UserID)
			ctx = utilctx.WithUserRole(ctx, accessToken.Claims.Custom.UserRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Authorize checks if user role in context is sufficient.
func Authorize(logger *zap.Logger, role domuser.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := utilctx.UserRoleFromCtx(r.Context())
			if !ok {
				httputil.WriteErrorResponse(
					r.Context(),
					util.WithCtx(r.Context(), logger),
					w,
					apierrors.NewUnauthorizedError(ErrMissingUserRole, "reading user role from ctx"),
				)
				return
			}

			if !userRole.IsSufficientToRole(role) {
				err := apierrors.NewForbiddenError(
					ErrInsufficientUserRole,
					"checking user role",
				).WithPublicMessage(ErrInsufficientUserRole.Error())
				httputil.WriteErrorResponse(r.Context(), util.WithCtx(r.Context(), logger), w, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseBearerToken(h http.Header) string {
	if h == nil {
		return ""
	}
	return strings.TrimPrefix(h.Get(authHeader), bearerSchema)
}
