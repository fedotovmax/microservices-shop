package errs

import (
	"errors"
	"fmt"
)

var ErrInvalidSession = errors.New("invalid session")

var ErrSessionNotFound = errors.New("session not found")

var ErrSessionExpired = errors.New("session is expired")

var ErrSessionRevoked = errors.New("session is revoked")

type UserSessionRevokedError struct {
	Email string
	UID   string
	SID   string
}

func NewUserSessionRevokedError(email, uid, sid string) *UserSessionRevokedError {
	return &UserSessionRevokedError{
		Email: email,
		UID:   uid,
		SID:   sid,
	}
}

func (err *UserSessionRevokedError) Error() string {
	return fmt.Sprintf("session revoked: sid=%s; uid=%s; email=%s", err.SID, err.UID, err.Email)
}

func (err *UserSessionRevokedError) Unwrap() error {
	return ErrSessionRevoked
}

var ErrUserSessionsInBlackList = errors.New("session in blacklist")

type UserSessionsInBlacklistError struct {
	Email string
	UID   string
}

func NewUserSessionsInBlacklistError(email, uid string) *UserSessionsInBlacklistError {
	return &UserSessionsInBlacklistError{
		Email: email,
		UID:   uid,
	}
}

func (err *UserSessionsInBlacklistError) Error() string {
	return fmt.Sprintf("sessions in blacklist: uid=%s; email=%s", err.UID, err.Email)
}

func (err *UserSessionsInBlacklistError) Unwrap() error {
	return ErrUserSessionsInBlackList
}

var ErrAgentLooksLikeBot = errors.New("the user agent looks like a bot")

var ErrInvalidSessionIP = errors.New("invalid session IP")
