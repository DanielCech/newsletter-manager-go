package graph

import (
	"context"
	"errors"
	"time"

	"strv-template-backend-go-api/api/graphql/middleware"
	apierrors "strv-template-backend-go-api/types/errors"

	netx "go.strv.io/net"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// ErrorPresenter modifies error response based on api error fields.
// Also logs the error with all extensions included.
func ErrorPresenter(logger *zap.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		result := graphql.DefaultErrorPresenter(ctx, err)
		if result.Unwrap() == nil {
			// An error directly from graphql, like a validation error.
			return result
		}
		if result.Extensions == nil {
			result.Extensions = make(map[string]any)
		}
		if requestID := netx.RequestIDFromCtx(ctx); requestID != "" {
			result.Extensions["requestId"] = requestID
		}

		isInternalError := false
		var apiErr *apierrors.Error
		if errors.As(err, &apiErr) {
			result.Message = apiErr.PublicMessage
			result.Extensions["code"] = apiErr.Code
			if apiErr.Code == apierrors.CodeUnknown {
				isInternalError = true
			}
		} else {
			// Unknown error type.
			result.Message = "internal server error"
			result.Extensions["code"] = apierrors.CodeUnknown
			isInternalError = true
		}

		l := logger.With(
			zap.String("path", result.Path.String()),
			zap.String("query", graphql.GetOperationContext(ctx).RawQuery),
			zap.Any("extensions", result.Extensions),
			zap.Error(err),
		)
		if requestStartTime, ok := middleware.StartTimeFromCtx(ctx); ok {
			l = l.With(zap.Int64("duration_ms", time.Since(requestStartTime).Milliseconds()))
		}

		if isInternalError {
			l.Error("request processed")
		} else {
			l.Info("request processed")
		}

		return result
	}
}
