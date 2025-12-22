package inputs

import (
	"errors"
	"regexp"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/validation"
)

var UpperLettersRegexp = regexp.MustCompile(`[A-Z]`)
var LowerLettersRegexp = regexp.MustCompile(`[a-z]`)
var DigitRegexp = regexp.MustCompile(`\d`)
var SpecialRegexp = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\;:'",.<>/?]`)
var NameRegexp = regexp.MustCompile(`^\p{L}+(?:[ '’]\p{L}+)*$`)
var PhoneRegexp = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func addValidationError(f, l, m string, err error) violations.ValidationError {
	return violations.ValidationError{
		Field:       f,
		Reason:      err.Error(),
		Description: "ValidationError",
		LocalizedMessage: &violations.LocalizedMessage{
			Locale:  l,
			Message: m,
		},
	}
}

func validateName(name string, locale string) (string, error) {

	err := validation.LengthRange(name, 1, 100)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationStrSymbolsRange, 1, 100)
		return msg, err
	}

	err = validation.Regex(name, NameRegexp)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationFullName)
		return msg, err
	}
	return "", nil
}

func validatePassword(password string) error {
	err := validation.Regex(password, UpperLettersRegexp)
	if err != nil {
		return err
	}

	err = validation.Regex(password, LowerLettersRegexp)

	if err != nil {
		return err
	}

	err = validation.Regex(password, DigitRegexp)

	if err != nil {
		return err
	}

	err = validation.Regex(password, SpecialRegexp)

	if err != nil {
		return err
	}

	return nil
}

func ValidateUUID(id string, locale string) error {

	_, err := validation.IsUUID(id)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationUUID)

		ve := violations.ValidationErrors{}
		ve = append(ve, addValidationError("ID", locale, msg, err))

		return ve
	}

	return nil

}

func validateGender(gender *domain.GenderValue) error {

	if gender == nil {
		return nil
	}

	if gender.IsValid() {
		return nil
	}

	return errors.New("invalid gender")

}
