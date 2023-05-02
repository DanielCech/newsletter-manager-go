package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	httputil "newsletter-manager-go/api/rest/util"
	apierrors "newsletter-manager-go/types/errors"
	"newsletter-manager-go/util"
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
	ParseAccessToken(data string) (string, error)
}

// Authenticate parses bearer token from authorization header.
// Custom claims parsed from access token are passed to context.
func Authenticate(logger *zap.Logger) func(http.Handler) http.Handler {
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

			///*accessToken*/
			//_, err := tokenParser.ParseAccessToken(token)
			//if err != nil {
			//	httputil.WriteErrorResponse(
			//		r.Context(),
			//		util.WithCtx(r.Context(), logger),
			//		w,
			//		apierrors.NewUnauthorizedError(err, "parsing access token"),
			//	)
			//	return
			//}

			ctx := context.Background()
			//ctx := utilctx.WithAuthorID(r.Context(), accessToken.Claims.AuthorID)
			//ctx = utilctx.WithUserRole(ctx, accessToken.Claims.Custom.UserRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return next
	}
}

//
//// Authorize checks if user role in context is sufficient.
//func Authorize(logger *zap.Logger, role domauthor.Role) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			userRole, ok := utilctx.UserRoleFromCtx(r.Context())
//			if !ok {
//				httputil.WriteErrorResponse(
//					r.Context(),
//					util.WithCtx(r.Context(), logger),
//					w,
//					apierrors.NewUnauthorizedError(ErrMissingUserRole, "reading user role from ctx"),
//				)
//				return
//			}
//
//			if !userRole.IsSufficientToRole(role) {
//				err := apierrors.NewForbiddenError(
//					ErrInsufficientUserRole,
//					"checking user role",
//				).WithPublicMessage(ErrInsufficientUserRole.Error())
//				httputil.WriteErrorResponse(r.Context(), util.WithCtx(r.Context(), logger), w, err)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}

func parseBearerToken(h http.Header) string {
	//if h == nil {
	//	return ""
	//}
	//return strings.TrimPrefix(h.Get(authHeader), bearerSchema)

	return ""
}
