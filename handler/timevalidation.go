package handler

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	//en_translations "github.com/go-playground/validator/v10/translations/en"
)



func Myvalidate(f1 validator.FieldLevel) bool {
	fieldvalue := f1.Field().Int()
	return fieldvalue == 10

}

func HourValidate(f1 validator.FieldLevel) bool {
	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	return timeRegex.MatchString(f1.Field().String())
}

func HourSecondValidate(f1 validator.FieldLevel) bool {
	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	return timeRegex.MatchString(f1.Field().String())
}
