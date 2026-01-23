package errs

import (
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrBadCredentials = errors.New(keys.UserBadCredentials)
var ErrEmailNotVerified = errors.New(keys.UserEmailNotVerified)
var ErrUserDeleted = errors.New(keys.UserDeleted)

type UserDeletedError struct {
	ErrCode           string
	DeletedAt         time.Time
	LastChanceRestore time.Time
}

func NewUserDeletedError(code string, deletedAt, lastChance time.Time) *UserDeletedError {
	return &UserDeletedError{
		ErrCode:           code,
		DeletedAt:         deletedAt,
		LastChanceRestore: lastChance,
	}
}

func (err *UserDeletedError) Error() string {
	return fmt.Sprintf("error code: %s; deleted at: %s; last chance to restore account: %s", err.ErrCode, err.DeletedAt, err.LastChanceRestore)
}

func (err *UserDeletedError) Unwrap() error {
	return ErrUserDeleted
}
