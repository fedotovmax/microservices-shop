package inputs

import "github.com/fedotovmax/grpcutils/violations"

type UUID struct {
	uuid string
}

func NewUUIDInput() *UUID {
	return &UUID{}
}

func (i *UUID) SetUUID(uuid string) {
	i.uuid = uuid
}

func (i *UUID) GetUUID() string {
	return i.uuid
}

func (i *UUID) Validate(locale string, fieldName string) error {
	var validationErrors violations.ValidationErrors

	msg, err := validateUUID(i.uuid, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError(fieldName, locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil

}
