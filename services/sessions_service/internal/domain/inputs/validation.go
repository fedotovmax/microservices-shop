package inputs

import (
	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/validation"
)

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

func validateUUID(uuid string, locale string) (string, error) {
	_, err := validation.IsUUID(uuid)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationUUID)
		return msg, err
	}
	return "", nil
}

func validateIP(ip string, locale string) (string, error) {
	_, err := validation.IsIPV4(ip)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationIP)
		return msg, err
	}
	return "", nil
}

func validateEmptyString(value string, locale string) (string, error) {
	err := validation.MinLength(value, 1)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationStrSymbolsMin, 1)
		return msg, err
	}
	return "", nil

}
