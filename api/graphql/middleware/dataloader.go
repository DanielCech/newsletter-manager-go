package middleware

import (
	"net/http"

	"strv-template-backend-go-api/database/sql"
	userdataloader "strv-template-backend-go-api/domain/user/postgres/dataloader"
)

// DataLoader passes data loaders to context so they are available in resolvers.
func DataLoader(dataSource sql.DataSource) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userLoader := userdataloader.New(dataSource)
			ctx := userdataloader.WithLoader(r.Context(), userLoader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
