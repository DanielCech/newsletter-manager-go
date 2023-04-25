package errors

// Code is useful for converting to HTTP status code but can be also used in Graphql
// as an additional information about the error.
// In general it's a value which is machine readable.
type Code string

const (
	CodeBadRequest      = Code("ERR_BAD_REQUEST")
	CodeUnauthorized    = Code("ERR_UNAUTHORIZED")
	CodeForbidden       = Code("ERR_FORBIDDEN")
	CodeNotFound        = Code("ERR_NOT_FOUND")
	CodeAlreadyExists   = Code("ERR_ALREADY_EXISTS")
	CodeExpired         = Code("ERR_EXPIRED")
	CodePayloadTooLarge = Code("ERR_PAYLOAD_TOO_LARGE")
	CodeInvalidBody     = Code("ERR_INVALID_BODY")
	CodeUnknown         = Code("ERR_UNKNOWN")
)
