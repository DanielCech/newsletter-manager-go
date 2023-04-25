package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	Validate = validator.New()
)

func init() {
	const tagName = "json"
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get(tagName), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
