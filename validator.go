package pkgep

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func InitValidator() {
	Validator = validator.New()

	// custom tag for triming out empty string for string type data
	Validator.RegisterValidation("trimmed", func(fl validator.FieldLevel) bool {
		if field, ok := fl.Field().Interface().(string); ok {
			trimmed := strings.TrimSpace(field)
			fl.Field().SetString(trimmed)
			return true
		}
		return false
	})
}

func ExtractValidationErrors(err error) string {
	if _, ok := err.(validator.ValidationErrors); ok {
		errs := err.(validator.ValidationErrors)
		errorSlc := []string{}
		errorMap := make(map[string]bool)

		for _, e := range errs {
			if errorMap[e.Field()] {
				continue
			}
			// turn first index of string to lowercase.
			runes := []rune(e.Field())
			runes[0] = unicode.ToLower(runes[0])
			errorSlc = append(errorSlc, string(runes))
			errorMap[e.Field()] = true
		}
		return strings.Join(errorSlc, ",")
	}
	return ""
}

// 👉 extract a validation error
func FormatValidateError(errs validator.ValidationErrors) string {
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			return fmt.Sprintf("The field '%s' is required", err.Field())
		case "max":
			return fmt.Sprintf("The field '%s' must not exceed %s characters", err.Field(), err.Param())
		case "url":
			return fmt.Sprintf("The field '%s' must be a valid URL", err.Field())
		}
	}
	return "Invalid input"
}
