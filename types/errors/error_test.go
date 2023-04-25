package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errTest = errors.New("test error")
)

func Test_NewError(t *testing.T) {
	message := "creating session"
	apiError := NewError(errTest, message, CodeUnknown)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Empty(t, apiError.PublicMessage)
	assert.Equal(t, apiError.Code, CodeUnknown)
	assert.Empty(t, apiError.Data)
}

func Test_Error_WithPublicMessage(t *testing.T) {
	message := "creating session"
	apiError := &Error{}
	apiError = apiError.WithPublicMessage(message)
	assert.Equal(t, message, apiError.PublicMessage)
}

func Test_Error_WithData(t *testing.T) {
	data := map[string]any{"id": 1}
	apiError := &Error{}
	apiError = apiError.WithData(data)
	assert.Equal(t, data, apiError.Data)
}

func Test_Error_Error(t *testing.T) {
	apiError := &Error{}
	assert.Empty(t, apiError.Error())

	apiError = &Error{err: errTest}
	assert.EqualError(t, apiError, errTest.Error())

	message := "creating session"
	apiError = &Error{Message: message}
	assert.Equal(t, message, apiError.Error())

	apiError = &Error{
		err:     errTest,
		Message: message,
	}
	assert.EqualError(t, apiError, fmt.Sprintf("%s: %s", message, errTest.Error()))
}

func Test_Error_Unwrap(t *testing.T) {
	apiError := &Error{
		err: errTest,
	}
	assert.Equal(t, errTest, apiError.Unwrap())
}

func Test_NewBadRequestError(t *testing.T) {
	message := "creating session"
	apiError := NewBadRequestError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeBadRequest, apiError.Code)
}

func Test_NewUnauthorizedError(t *testing.T) {
	message := "creating session"
	apiError := NewUnauthorizedError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeUnauthorized, apiError.Code)
}

func Test_NewForbiddenError(t *testing.T) {
	message := "creating session"
	apiError := NewForbiddenError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeForbidden, apiError.Code)
}

func Test_NewNotFoundError(t *testing.T) {
	message := "creating session"
	apiError := NewNotFoundError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeNotFound, apiError.Code)
}

func Test_NewAlreadyExistsError(t *testing.T) {
	message := "creating session"
	apiError := NewAlreadyExistsError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeAlreadyExists, apiError.Code)
}

func Test_NewExpiredError(t *testing.T) {
	message := "creating session"
	apiError := NewExpiredError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeExpired, apiError.Code)
}

func Test_NewPayloadTooLargeError(t *testing.T) {
	message := "creating session"
	apiError := NewPayloadTooLargeError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodePayloadTooLarge, apiError.Code)
}

func Test_NewInvalidBodyError(t *testing.T) {
	message := "creating session"
	apiError := NewInvalidBodyError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeInvalidBody, apiError.Code)
}

func Test_NewUnknownError(t *testing.T) {
	message := "creating session"
	apiError := NewUnknownError(errTest, message)
	assert.Equal(t, errTest, apiError.err)
	assert.Equal(t, message, apiError.Message)
	assert.Equal(t, CodeUnknown, apiError.Code)
}
