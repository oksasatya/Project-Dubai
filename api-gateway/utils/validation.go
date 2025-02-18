package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func FormatValidationError(model interface{}, errors validator.ValidationErrors) string {
	var errorMessages []string
	modelType := reflect.TypeOf(model).Elem()

	for _, err := range errors {
		// Get the field name from the models's struct
		field, found := modelType.FieldByName(err.Field())
		fieldName := err.Field()
		if found {
			fieldName = field.Tag.Get("json")
			if fieldName == "" {
				fieldName = err.Field()
			}
		}

		// Generate custom error message based on the validation tag
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "numeric":
			message = fmt.Sprintf("%s must be a number", fieldName)
		case "unique":
			message = fmt.Sprintf("%s already exists", fieldName)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", fieldName, err.Param())
		case "max":
			message = fmt.Sprintf("%s must not be longer than %s characters", fieldName, err.Param())
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, err.Param())
		case "gt":
			message = fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", fieldName, err.Param())
		case "lt":
			message = fmt.Sprintf("%s must be less than %s", fieldName, err.Param())
		case "eqfield":
			message = fmt.Sprintf("%s must be equal to %s", fieldName, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}

		errorMessages = append(errorMessages, message)
	}

	return strings.Join(errorMessages, ", ")
}
