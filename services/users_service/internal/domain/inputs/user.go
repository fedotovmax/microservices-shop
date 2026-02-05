package inputs

import (
	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

type SessionAction struct {
	*CreateUser
}

func NewSessionAction() *SessionAction {
	return &SessionAction{
		CreateUser: NewCreateUser(),
	}
}

type CreateUser struct {
	email    string
	password string
}

func NewCreateUser() *CreateUser {

	return &CreateUser{}
}

func (i *CreateUser) Validate(locale string) error {
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

func (i *CreateUser) GetEmail() string {
	return i.email
}

func (i *CreateUser) GetPassword() string {
	return i.password
}

func (i *CreateUser) SetEmail(email string) {
	i.email = email
}

func (i *CreateUser) SetPassword(password string) {
	i.password = password
}

type UpdateUser struct {
	userID     string
	birthDate  *string
	lastName   *string
	firstName  *string
	middleName *string
	avatarURL  *string
	gender     *domain.GenderValue
}

func NewUpdateUser() *UpdateUser {
	return &UpdateUser{}
}

func (i *UpdateUser) Validate(locale string) error {

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

func (i *UpdateUser) GetBirthDate() *string {
	return i.birthDate
}

func (i *UpdateUser) GetFirstName() *string {
	return i.firstName
}

func (i *UpdateUser) GetLastName() *string {
	return i.lastName
}

func (i *UpdateUser) GetMiddleName() *string {
	return i.middleName
}

func (i *UpdateUser) GetAvatarURL() *string {
	return i.avatarURL
}

func (i *UpdateUser) GetGender() *domain.GenderValue {
	return i.gender
}

func (i *UpdateUser) GetUserID() string {
	return i.userID
}

func (i *UpdateUser) SetFromProto(req *userspb.UpdateUserProfileRequest) {

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

func (i *UpdateUser) SetBirthDate(b *string) {
	i.birthDate = b
}

func (i *UpdateUser) SetFirstName(f *string) {
	i.firstName = f
}

func (i *UpdateUser) SetLastName(l *string) {
	i.lastName = l
}

func (i *UpdateUser) SetMiddleName(m *string) {
	i.middleName = m
}

func (i *UpdateUser) SetAvatarURL(url *string) {
	i.avatarURL = url
}

func (i *UpdateUser) SetGender(g *domain.GenderValue) {
	i.gender = g
}
