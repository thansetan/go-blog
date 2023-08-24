package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type InputError struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

func ValidationError(valErr error) any {
	var ve validator.ValidationErrors
	if errors.As(valErr, &ve) {
		out := make([]InputError, len(ve))
		for i, fe := range ve {
			out[i] = InputError{strings.ToLower(fe.Field()), errorMessage(fe.Tag(), fe.Param())}
		}

		return out
	}
	return valErr
}

func errorMessage(tag, param string) string {
	fmt.Println(tag)
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "This field has to be a valid email"
	case "min":
		return fmt.Sprintf("This field has to be at least %s characters long", param)
	case "max":
		return fmt.Sprintf("This field length can't exceed %s characters", param)
	case "password":
		return "This field must contains at least 1 numerical character, 1 symbol and 1 uppercase letter"
	default:
		return ""
	}
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasNumeric := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)

	return hasNumeric && hasSymbol && hasUppercase
}
