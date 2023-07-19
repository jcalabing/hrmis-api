package functions

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateNonEmpty(fl validator.FieldLevel) bool {
	field := fl.Field()
	value := strings.TrimSpace(field.String())
	return value != ""
}
