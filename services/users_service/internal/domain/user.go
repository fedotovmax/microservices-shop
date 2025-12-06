package domain

import (
	"time"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/validation"
)

type CreateUserInput struct {
	email    string
	password string
}

func NewCreateUserInput(email, password string) *CreateUserInput {

	return &CreateUserInput{email: email, password: password}
}

func (i CreateUserInput) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	_, err := validation.IsEmail(i.email)

	if err != nil {

		msg, lerr := i18n.Local.Get(locale, keys.ValidationEmail)

		if lerr != nil {
			return lerr
		}

		validationErrors = append(validationErrors, AddValidationError("Email", locale, msg, err))
	}

	err = validatePassword(i.password)

	if err != nil {

		msg, lerr := i18n.Local.Get(locale, keys.ValidationPassword)

		if lerr != nil {
			return lerr
		}

		validationErrors = append(validationErrors,
			AddValidationError("Password", locale, msg, err))
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

func (i *CreateUserInput) SetPasswordHash(hash string) {
	i.password = hash
}

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Profile      Profile
	ID           string
	Email        string
	Phone        *string
	PasswordHash string
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Profile struct {
	UpdatedAt  time.Time
	BirthDate  *time.Time
	LastName   *string
	FirstName  *string
	MiddleName *string
	AvatarURL  *string
	Gender     *Gender
}

type UserCreatedEvent struct {
	ID string
}
