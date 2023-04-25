package errors

import (
	"fmt"
)

// Error wraps an existing error with additional information.
type Error struct {
	err           error
	Message       string
	PublicMessage string
	Code          Code
	Data          any
}

// NewError returns new instance of Error.
func NewError(err error, message string, code Code) *Error {
	return &Error{
		err:     err,
		Message: message,
		Code:    code,
		Data:    nil,
	}
}

// WithPublicMessage adds public message into the error.
func (e *Error) WithPublicMessage(msg string) *Error {
	e.PublicMessage = msg
	return e
}

// WithData adds public data into the error.
func (e *Error) WithData(data any) *Error {
	e.Data = data
	return e
}

// Error returns error message.
func (e *Error) Error() string {
	if e.Message != "" {
		if e.err != nil {
			return fmt.Sprintf("%s: %s", e.Message, e.err.Error())
		}
		return e.Message
	}
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// Unwrap returns inner error.
func (e *Error) Unwrap() error {
	return e.err
}

// NewBadRequestError returns new instance of Error with CodeBadRequest code.
func NewBadRequestError(err error, message string) *Error {
	return NewError(err, message, CodeBadRequest)
}

// NewUnauthorizedError returns new instance of Error with CodeUnauthorized code.
func NewUnauthorizedError(err error, message string) *Error {
	return NewError(err, message, CodeUnauthorized)
}

// NewForbiddenError returns new instance of Error with CodeForbidden code.
func NewForbiddenError(err error, message string) *Error {
	return NewError(err, message, CodeForbidden)
}

// NewNotFoundError returns new instance of Error with CodeNotFound code.
func NewNotFoundError(err error, message string) *Error {
	return NewError(err, message, CodeNotFound)
}

// NewAlreadyExistsError returns new instance of Error with CodeAlreadyExists code.
func NewAlreadyExistsError(err error, message string) *Error {
	return NewError(err, message, CodeAlreadyExists)
}

// NewExpiredError returns new instance of Error with CodeExpired code.
func NewExpiredError(err error, message string) *Error {
	return NewError(err, message, CodeExpired)
}

// NewPayloadTooLargeError returns new instance of Error with CodePayloadTooLarge code.
func NewPayloadTooLargeError(err error, message string) *Error {
	return NewError(err, message, CodePayloadTooLarge)
}

// NewInvalidBodyError returns new instance of Error with CodeInvalidBody code.
func NewInvalidBodyError(err error, message string) *Error {
	return NewError(err, message, CodeInvalidBody)
}

// NewUnknownError returns new instance of Error with CodeUnknown code.
func NewUnknownError(err error, message string) *Error {
	return NewError(err, message, CodeUnknown)
}
