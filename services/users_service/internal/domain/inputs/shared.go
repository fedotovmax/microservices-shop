package inputs

import "github.com/fedotovmax/grpcutils/violations"

type UUIDInput struct {
	uuid string
}

func NewUUIDInput() *UUIDInput {
	return &UUIDInput{}
}

func (i *UUIDInput) SetUUID(uuid string) {
	i.uuid = uuid
}

func (i *UUIDInput) GetUUID() string {
	return i.uuid
}

func (i *UUIDInput) Validate(locale string, fieldName string) error {
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
