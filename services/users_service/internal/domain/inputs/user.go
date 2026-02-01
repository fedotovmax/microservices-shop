package inputs

import (
	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
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

	msg, err := validateEmail(i.email, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("Email", locale, msg, err))
	}

	msg, err = validatePassword(i.password, locale)

	if err != nil {
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
	userID     string
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

	msg, err := validateUUID(i.userID, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("UserID", locale, msg, err))
	}

	msg, err = validateGender(i.gender, locale)

	if err != nil {
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
		msg, err := validateFilePath(*i.avatarURL, locale)
		if err != nil {
			validationErrors = append(validationErrors, addValidationError("AvatarURL", locale, msg, err))
		}
	}

	if i.birthDate != nil {
		msg, err := validateDateString(*i.birthDate, keys.DateFormat, locale)
		if err != nil {
			validationErrors = append(validationErrors, addValidationError("BirthDate", locale, msg, err))
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

func (i *UpdateUserInput) GetUserID() string {
	return i.userID
}

func (i *UpdateUserInput) SetFromProto(req *userspb.UpdateUserProfileRequest) {

	if req != nil && req.Data != nil {
		i.avatarURL = req.Data.AvatarUrl
		i.birthDate = req.Data.BirthDate
		i.firstName = req.Data.FirstName
		i.lastName = req.Data.LastName
		i.middleName = req.Data.MiddleName
		i.gender = domain.GenderFromProto(req.Data.GenderValue)
		i.userID = req.UserId
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

type FindUserByIDInput struct {
	userID string
}

func NewFindUserByIDInput() *FindUserByIDInput {
	return &FindUserByIDInput{}
}

func (i *FindUserByIDInput) SetUserID(id string) {
	i.userID = id
}

func (i *FindUserByIDInput) GetUserID() string {
	return i.userID
}

func (i *FindUserByIDInput) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	msg, err := validateUUID(i.userID, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("UserID", locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil

}
