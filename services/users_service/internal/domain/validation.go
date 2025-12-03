package domain

import (
	"regexp"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/validation"
)

var UpperLettersRegexp = regexp.MustCompile(`[A-Z]`)
var LowerLettersRegexp = regexp.MustCompile(`[a-z]`)
var DigitRegexp = regexp.MustCompile(`\d`)
var SpecialRegexp = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\;:'",.<>/?]`)

var PhoneRegexp = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func AddValidationError(f, l, m string, err error) violations.ValidationError {
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
