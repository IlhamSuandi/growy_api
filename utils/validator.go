package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type ValidationError struct {
	Errors map[string]string
}

func (ve ValidationError) Error() string {
	var errMsgs []string
	for field, msg := range ve.Errors {
		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(errMsgs, ", ")
}

func getErrorMessage(err validator.FieldError) string {
	errorList := map[string]string{
		"required":              "field is required",
		"email":                 "invalid email address",
		"min":                   fmt.Sprintf("%s is too short, minimum length is %s", err.Field(), err.Param()),
		"max":                   fmt.Sprintf("%s is too long, maximum length is %s", err.Field(), err.Param()),
		"alphanum":              fmt.Sprintf("%s only allow alphabet and number", err.Field()),
		"contains_uppercase":    fmt.Sprintf("%s must have uppercase letter", err.Field()),
		"contains_num":          fmt.Sprintf("%s must have number", err.Field()),
		"contains_special_char": fmt.Sprintf("%s must have special character", err.Field()),
		"exclude_special_char":  fmt.Sprintf("%s only allow alphabet and number", err.Field()),
	}

	if errorMessage, ok := errorList[err.Tag()]; ok {
		return errorMessage
	}

	return err.Tag()
}

func regexMatch(fl validator.FieldLevel, pattern string) bool {
	field := fl.Field().String()
	re := regexp.MustCompile(pattern)
	return re.MatchString(field)
}

func RegisterValidator() {
	// register contain uppercase
	Validate.RegisterValidation("contains_uppercase", func(fl validator.FieldLevel) bool {
		uppercasePattern := `[A-Z]`
		return regexMatch(fl, uppercasePattern)
	})

	// register special char
	Validate.RegisterValidation("contains_special_char", func(fl validator.FieldLevel) bool {
		specialCharPattern := `[!@#\$%\^&\*\(\)\-_=\+\[\]\{\};:'"<>,\.\?\/\\|~]`
		return regexMatch(fl, specialCharPattern)
	})

	// register contains number
	Validate.RegisterValidation("contains_num", func(fl validator.FieldLevel) bool {
		numPattern := `[0-9]`
		return regexMatch(fl, numPattern)
	})

	// only allow alphabet and number
	Validate.RegisterValidation("exclude_special_char", func(fl validator.FieldLevel) bool {
		numPattern := `^[a-zA-Z0-9 ]+$`
		return regexMatch(fl, numPattern)
	})
}

func ValidateStruct(s interface{}) error {
	errMap := make(map[string]string)
	if err := Validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			error_message := getErrorMessage(err)
			errMap[err.Field()] = error_message
		}
		return ValidationError{Errors: errMap}
	}
	return nil
}
