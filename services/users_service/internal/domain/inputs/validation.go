package inputs

import (
	"errors"
	"regexp"
	"time"

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

func validateEmail(email string, locale string) (string, error) {
	_, err := validation.IsEmail(email)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationEmail)
		return msg, nil
	}

	return "", nil

}

func validateUUID(uuid string, locale string) (string, error) {
	_, err := validation.IsUUID(uuid)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationUUID)
		return msg, err
	}
	return "", nil
}

func validateDateString(date string, format string, locale string) (string, error) {
	_, err := time.Parse(format, date)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationDateFormat)
		return msg, err
	}
	return "", nil
}

func validateFilePath(path string, locale string) (string, error) {
	err := validation.IsFilePath(path)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationStrFilePath)
		return msg, err
	}
	return "", nil
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

func validatePassword(password string, locale string) (string, error) {

	err := validation.MinLength(password, 8)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationStrSymbolsMin, 8)
		return msg, err
	}

	err = validation.Regex(password, UpperLettersRegexp)
	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationPassword)
		return msg, err
	}

	err = validation.Regex(password, LowerLettersRegexp)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationPassword)
		return msg, err
	}

	err = validation.Regex(password, DigitRegexp)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationPassword)
		return msg, err
	}

	err = validation.Regex(password, SpecialRegexp)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationPassword)
		return msg, err
	}

	return "", nil
}

func validateGender(gender *domain.GenderValue, locale string) (string, error) {

	if gender == nil {
		return "", nil
	}

	if gender.IsValid() {
		return "", nil
	}

	msg, _ := i18n.Local.Get(locale, keys.ValidationGender)

	return msg, errors.New("invalid gender")

}
