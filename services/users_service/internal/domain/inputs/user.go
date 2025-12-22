package inputs

import (
	"time"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/validation"
)

type SessionActionInput struct {
	*CreateUserInput
}

func NewSessionActionInput() *SessionActionInput {
	return &SessionActionInput{
		CreateUserInput: NewCreateUserInput(),
	}
}

type CreateUserInput struct {
	email    string
	password string
}

func NewCreateUserInput() *CreateUserInput {

	return &CreateUserInput{}
}

func (i *CreateUserInput) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	_, err := validation.IsEmail(i.email)

	if err != nil {

		msg, _ := i18n.Local.Get(locale, keys.ValidationEmail)

		validationErrors = append(validationErrors, addValidationError("Email", locale, msg, err))
	}

	err = validatePassword(i.password)

	if err != nil {

		msg, _ := i18n.Local.Get(locale, keys.ValidationPassword)

		validationErrors = append(validationErrors,
			addValidationError("Password", locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (i *CreateUserInput) GetEmail() string {
	return i.email
}

func (i *CreateUserInput) GetPassword() string {
	return i.password
}

func (i *CreateUserInput) SetEmail(email string) {
	i.email = email
}

func (i *CreateUserInput) SetPassword(password string) {
	i.password = password
}

type UpdateUserInput struct {
	birthDate  *string
	lastName   *string
	firstName  *string
	middleName *string
	avatarURL  *string
	gender     *domain.GenderValue
}

func NewUpdateUserInput() *UpdateUserInput {
	return &UpdateUserInput{}
}

func (i *UpdateUserInput) Validate(locale string) error {

	var validationErrors violations.ValidationErrors

	err := validateGender(i.gender)

	if err != nil {
		msg, _ := i18n.Local.Get(locale, keys.ValidationGender)
		validationErrors = append(validationErrors, addValidationError("Gender", locale, msg, err))
	}

	if i.firstName != nil {
		msg, err := validateName(*i.firstName, locale)
		if err != nil {
			validationErrors = append(validationErrors, addValidationError("FirstName", locale, msg, err))
		}
	}

	if i.lastName != nil {
		msg, err := validateName(*i.lastName, locale)
		if err != nil {
			validationErrors = append(validationErrors, addValidationError("LastName", locale, msg, err))
		}
	}

	if i.middleName != nil {
		msg, err := validateName(*i.middleName, locale)
		if err != nil {
			validationErrors = append(validationErrors, addValidationError("MiddleName", locale, msg, err))
		}
	}

	if i.avatarURL != nil {
		err := validation.IsFilePath(*i.avatarURL)
		if err != nil {
			msg, _ := i18n.Local.Get(locale, keys.ValidationStrFilePath)
			validationErrors = append(validationErrors, addValidationError("AvatarURL", locale, msg, err))
		}
	}

	if i.birthDate != nil {
		_, err := time.Parse(keys.DateFormat, *i.birthDate)
		if err != nil {
			msg, _ := i18n.Local.Get(locale, keys.ValidationDateFormat)
			validationErrors = append(validationErrors, addValidationError("AvatarURL", locale, msg, err))
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (i *UpdateUserInput) GetBirthDate() *string {
	return i.birthDate
}

func (i *UpdateUserInput) GetFirstName() *string {
	return i.firstName
}

func (i *UpdateUserInput) GetLastName() *string {
	return i.lastName
}

func (i *UpdateUserInput) GetMiddleName() *string {
	return i.middleName
}

func (i *UpdateUserInput) GetAvatarURL() *string {
	return i.avatarURL
}

func (i *UpdateUserInput) GetGender() *domain.GenderValue {
	return i.gender
}

func (i *UpdateUserInput) SetFromProto(req *userspb.UpdateUserProfileRequest) {

	if req != nil {
		i.avatarURL = req.AvatarUrl
		i.birthDate = req.BirthDate
		i.firstName = req.FirstName
		i.lastName = req.LastName
		i.middleName = req.MiddleName
		i.gender = domain.GenderFromProto(req.GenderValue)
	}
}

func (i *UpdateUserInput) SetBirthDate(b *string) {
	i.birthDate = b
}

func (i *UpdateUserInput) SetFirstName(f *string) {
	i.firstName = f
}

func (i *UpdateUserInput) SetLastName(l *string) {
	i.lastName = l
}

func (i *UpdateUserInput) SetMiddleName(m *string) {
	i.middleName = m
}

func (i *UpdateUserInput) SetAvatarURL(url *string) {
	i.avatarURL = url
}

func (i *UpdateUserInput) SetGender(g *domain.GenderValue) {
	i.gender = g
}
