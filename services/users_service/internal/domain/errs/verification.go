package errs

import "errors"

var ErrEmailNotVerified = errors.New("email not verified")
var ErrVerifyEmailLinkNotFound = errors.New("verify email link not found")
var ErrVerifyEmailLinkExpired = errors.New("verify email link expired")
var ErrUserEmailAlreadyVerified = errors.New("user email already verified!")

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
