package platform

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func initValidators(v *validator.Validate, trans ut.Translator) {
	en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterValidation("username", validateUsername)

	v.RegisterTranslation("username", trans, func(ut ut.Translator) error {
		return ut.Add("username", "{0} must start with a letter and contain only letters, numbers, and underscores", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username", fe.Field())
		return t
	})
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]*$", username)
	return matched
}
