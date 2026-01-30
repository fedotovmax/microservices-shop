package errs

import (
	"errors"
	"fmt"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserEmailAlreadyVerified = errors.New("user email already verified!")

var ErrVerifyEmailLinkNotFound = errors.New("verify email link not found")
var ErrVerifyEmailLinkExpired = errors.New("verify email link expired")

var ErrBadCredentials = errors.New("bad credentials")
var ErrEmailNotVerified = errors.New("email not verified")
var ErrUserDeleted = errors.New("user deleted")

type EmailNotVerifiedError struct {
	ErrCode string
	UID     string
}

func NewEmailNotVerifiedErrorError(code, uid string) *EmailNotVerifiedError {
	return &EmailNotVerifiedError{
		ErrCode: code,
		UID:     uid,
	}
}

func (err *EmailNotVerifiedError) Error() string {
	return "email not verified"
}

func (err *EmailNotVerifiedError) Unwrap() error {
	return ErrEmailNotVerified
}

type VerifyEmailLinkExpiredError struct {
	ErrCode string
	UID     string
}

func NewVerifyEmailLinkExpiredError(code, uid string) *VerifyEmailLinkExpiredError {
	return &VerifyEmailLinkExpiredError{
		ErrCode: code,
		UID:     uid,
	}
}

func (err *VerifyEmailLinkExpiredError) Error() string {
	return "verify email link is expired"
}

func (err *VerifyEmailLinkExpiredError) Unwrap() error {
	return ErrVerifyEmailLinkExpired
}

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
