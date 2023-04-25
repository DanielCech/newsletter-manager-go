package util

import (
	"context"
	"errors"
	"net/http"

	apierrors "strv-template-backend-go-api/types/errors"

	netx "go.strv.io/net"
	httpx "go.strv.io/net/http"

	"go.uber.org/zap"
)

const (
	defaultErrorStatusCode = http.StatusInternalServerError
)

// WriteResponse is helper function for writing HTTP response with status code and options.
func WriteResponse(logger *zap.Logger, w http.ResponseWriter, data any, statusCode int, opts ...httpx.ResponseOption) {
	if err := httpx.WriteResponse(w, data, statusCode, opts...); err != nil {
		logger.With(
			zap.Int("status_code", statusCode),
			zap.Any("data", data),
		).Error("writing http response", zap.Error(err))
	}
}

// WriteErrorResponse is helper function for writing HTTP error response.
func WriteErrorResponse(ctx context.Context, logger *zap.Logger, w http.ResponseWriter, err error) {
	statusCode, opts := ToStatusCodeWithOptions(err)
	opts = append(opts, httpx.WithRequestID(netx.RequestIDFromCtx(ctx)))

	if err := httpx.WriteErrorResponse(w, statusCode, opts...); err != nil {
		logger.With(
			zap.Int("status_code", statusCode),
		).Error("writing http error response", zap.Error(err))
	}
}

// ToStatusCodeWithOptions attempts to cast error to *apierrors.Error and sets HTTP status code based on
// apierror code. If error does not match, default http 500 and "ERR_UNKNOWN" are used
// Options which are always set with *apierrors.Error:
//   - httpx.WithError is for logging purposes.
//   - httpx.WithErrorCode sets code in JSON response.
//   - httpx.WithErrorMessage sets public message in JSON response.
//   - httpx.WithErrorData sets arbitrary public data in JSON response.
func ToStatusCodeWithOptions(err error) (statusCode int, opts []httpx.ErrorResponseOption) {
	var apiErr *apierrors.Error
	if !errors.As(err, &apiErr) {
		return defaultErrorStatusCode, []httpx.ErrorResponseOption{
			httpx.WithError(err),
			httpx.WithErrorCode(string(apierrors.CodeUnknown)),
		}
	}

	opts = []httpx.ErrorResponseOption{
		httpx.WithError(err),
		httpx.WithErrorCode(string(apiErr.Code)),
		httpx.WithErrorMessage(apiErr.PublicMessage),
		httpx.WithErrorData(apiErr.Data),
	}

	switch c := apiErr.Code; c {
	case apierrors.CodeBadRequest:
		statusCode = http.StatusBadRequest
	case apierrors.CodeUnauthorized:
		statusCode = http.StatusUnauthorized
	case apierrors.CodeForbidden:
		statusCode = http.StatusForbidden
	case apierrors.CodeNotFound:
		statusCode = http.StatusNotFound
	case apierrors.CodeAlreadyExists:
		statusCode = http.StatusConflict
	case apierrors.CodeExpired:
		statusCode = http.StatusGone
	case apierrors.CodePayloadTooLarge:
		statusCode = http.StatusRequestEntityTooLarge
	case apierrors.CodeInvalidBody:
		statusCode = http.StatusUnprocessableEntity
	case apierrors.CodeUnknown:
		fallthrough
	default:
		statusCode = defaultErrorStatusCode
	}

	return statusCode, opts
}
