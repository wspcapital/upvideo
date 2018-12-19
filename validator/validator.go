package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"strings"
)

var (
	validate *validator.Validate // Validator instance

	titleRegexp       = regexp.MustCompile("^[\\p{L}\\d_\\s\\-]*$")
	textRegexp        = regexp.MustCompile("^[\\p{L}\\d_~!@#$%^&*()`\\[\\]{};':,./<>?|\\s\\-]*$")
	confirmCodeRegexp = regexp.MustCompile("^[a-zA-Z0-9/_\\-]*$")
)

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
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

func JsonErrors(errs interface{}) (json map[string]interface{}) {
	json = make(map[string]interface{})
	if len(errs.(validator.ValidationErrors)) > 0 {
		for _, err := range errs.(validator.ValidationErrors) {
			json[err.Field()] = err.Value()
		}
	}
	return
}
