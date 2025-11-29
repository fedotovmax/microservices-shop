package domain

import (
	"regexp"
	"time"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/validation"
)

type CreateUserInput struct {
	email    string
	password string
}

var UpperLettersRegexp = regexp.MustCompile(`[A-Z]`)
var LowerLettersRegexp = regexp.MustCompile(`[a-z]`)
var DigitRegexp = regexp.MustCompile(`\d`)
var SpecialRegexp = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\;:'",.<>/?]`)

var PhoneRegexp = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

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

func NewCreateUserInput(email, password string) (CreateUserInput, error) {

	var validationErrors violations.ValidationErrors

	_, err := validation.IsEmail(email)

	if err != nil {
		validationErrors = append(validationErrors, violations.ValidationError{
			Field:       "Email",
			Reason:      err.Error(),
			Description: "ValidationError",
			LocalizedMessage: &violations.LocalizedMessage{
				Locale:  "ru",
				Message: "Введённый адрес электронной почты не соответствует требуемому формату",
			},
		})
	}

	err = validatePassword(password)
	if err != nil {
		validationErrors = append(validationErrors, violations.ValidationError{
			Field:       "Password",
			Reason:      err.Error(),
			Description: "ValidationError",
			LocalizedMessage: &violations.LocalizedMessage{
				Locale:  "ru",
				Message: `Пароль должен содержать хотя бы одну заглавную букву (A–Z), хотя бы одну строчную букву (a–z), хотя бы одну цифру (0–9), хотя бы один специальный символ, а также иметь общую длину от 8 до 64 символов.`,
			},
		})
	}

	if len(validationErrors) > 0 {
		return CreateUserInput{}, validationErrors
	}

	return CreateUserInput{email: email, password: password}, nil
}

func (i CreateUserInput) GetEmail() string {
	return i.email
}

func (i CreateUserInput) GetPassword() string {
	return i.password
}

func (i *CreateUserInput) SetPasswordHash(hash string) {
	i.password = hash
}

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ID           string
	Email        string
	Phone        *string
	PasswordHash string
	Profile      Profile
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
