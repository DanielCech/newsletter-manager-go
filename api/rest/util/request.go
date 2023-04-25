package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	apierrors "strv-template-backend-go-api/types/errors"
	customvalidator "strv-template-backend-go-api/types/validator"

	httpparam "go.strv.io/net/http/param"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// ParseRequestBody decodes and validates a struct.
// This function expects the request body to be a JSON object and target to be a pointer to expected struct.
// If the request body is invalid, it returns an error.
func ParseRequestBody(r *http.Request, target any) error {
	if err := parseRequestBody(r, target); err != nil {
		return err
	}
	return validateRequestInput(target)
}

func parseRequestBody(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		var e *http.MaxBytesError
		if errors.As(err, &e) {
			return apierrors.NewPayloadTooLargeError(err, "decoding json")
		}
		const publicErrMsg = "invalid json body"
		return apierrors.NewBadRequestError(err, "decoding json").WithPublicMessage(publicErrMsg)
	}
	return nil
}

func validateRequestInput(target any) error {
	if err := customvalidator.Validate.Struct(target); err != nil {
		e := validator.ValidationErrors{}
		if !errors.As(err, &e) {
			return apierrors.NewInvalidBodyError(err, "validating content")
		}
		var fields []map[string]any
		for _, ve := range e {
			fieldErrorDescription := map[string]any{
				"name": ve.Field(),
			}
			fields = append(fields, fieldErrorDescription)
		}
		data := map[string]any{
			"invalidFields": fields,
		}
		return apierrors.NewInvalidBodyError(err, "validating content").WithData(data)
	}
	return nil
}

// ParseRequestInput is a signature.InputGetterFunc that fills the struct from all the following: path params, query params, request body.
// input value for signature.WrapHandler (and related) inner handlers.
// Although query and path parameters are not used in this project, this is a way how to parse values from those parameters
// using go.strv.io/net/http/param package.
func ParseRequestInput(r *http.Request, target any) error {
	// Don't call json unmarshal if target has no json tag, which means request body may be empty,
	// as only expected input are path and query parameters.
	// causes error if json.Unmarshal is called on empty body (EOF)
	if hasStructJSONTag(target) {
		if err := parseRequestBody(r, target); err != nil {
			return err
		}
	}

	// Is after parseRequestBody, as parseRequestBody possibly fills all fields,
	// even those not tagged with json, but tagged as query or path parameter.
	// This way, such filled fields will be reassigned in httpparam.ParamParser.
	if err := httpparam.DefaultParser().WithPathParamFunc(chi.URLParam).Parse(r, target); err != nil {
		return err
	}

	return validateRequestInput(target)
}

func hasStructJSONTag(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		_, exists := t.Field(i).Tag.Lookup("json")
		if exists {
			return true
		}
	}
	return false
}

type ID interface {
	UnmarshalText(data []byte) error
}

// GetPathID parses and returns ID path parameter based on UUID from id package.
func GetPathID[TID any, TPtrID interface {
	*TID
	ID
}](r *http.Request, paramName string) (pathID TID, err error) {
	i := TPtrID(new(TID))
	param := chi.URLParam(r, paramName)
	if err = i.UnmarshalText([]byte(param)); err != nil {
		publicErrMsg := fmt.Sprintf("invalid path id parameter %q", paramName)
		return pathID, apierrors.NewBadRequestError(err, "unmarshalling text").WithPublicMessage(publicErrMsg)
	}
	return *i, nil
}

// GetQueryID parses and returns ID query parameter based on UUID from id package.
func GetQueryID[TID any, TPtrID interface {
	*TID
	ID
}](r *http.Request, paramName string) (queryID TPtrID, err error) {
	queryID = TPtrID(new(TID))
	if !r.URL.Query().Has(paramName) {
		return queryID, nil
	}
	if err = queryID.UnmarshalText([]byte(r.URL.Query().Get(paramName))); err != nil {
		publicErrMsg := fmt.Sprintf("invalid query id parameter %q", paramName)
		return queryID, apierrors.NewBadRequestError(err, "unmarshalling text").WithPublicMessage(publicErrMsg)
	}
	return queryID, nil
}
