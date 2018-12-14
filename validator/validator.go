package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var (
	validate *validator.Validate // Validator instance

	titleRegexp       = regexp.MustCompile("^[\\p{L}\\d_\\s\\-]*$")
	textRegexp        = regexp.MustCompile("^[\\p{L}\\d_~!@#$%^&*()`\\[\\]{};':,./<>?|\\s\\-]*$")
	confirmCodeRegexp = regexp.MustCompile("^[a-zA-Z0-9/_\\-]*$")
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("title", validateTitle)
	_ = validate.RegisterValidation("text", validateText)
	_ = validate.RegisterValidation("confirm_code", validateConfirmCode)
}

func GetValidatorInstance() *validator.Validate {
	return validate
}

func validateTitle(fl validator.FieldLevel) bool {
	return titleRegexp.MatchString(fl.Field().String())
}

func validateText(fl validator.FieldLevel) bool {
	return textRegexp.MatchString(fl.Field().String())
}

func validateConfirmCode(fl validator.FieldLevel) bool {
	return confirmCodeRegexp.MatchString(fl.Field().String())
}
