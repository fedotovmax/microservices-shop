package errs

import (
	"errors"
	"fmt"
	"time"
)

var ErrTrustTokenNotFound = errors.New("trust token not found or expired or revoked")

var ErrUserDeleted = errors.New("user is deleted")

var ErrInvalidSession = errors.New("invalid session")

var ErrSessionNotFound = errors.New("session not found")

var ErrSessionExpired = errors.New("session is expired")

var ErrSessionRevoked = errors.New("session is revoked")

var ErrUserSessionsInBlackList = errors.New("user in blacklist")

var ErrBlacklistCodeExpired = errors.New("blacklist code expired")

var ErrBadBlacklistCode = errors.New("bad blacklist code")

var ErrAgentLooksLikeBot = errors.New("the user agent looks like a bot")

var ErrInvalidSessionIP = errors.New("invalid session IP")

var ErrLoginFromNewIPOrDevice = errors.New("login from new device")

var ErrBypassCodeExpired = errors.New("bypass code expired")

var ErrBadBypassCode = errors.New("invalid bypass code")

type UserSessionsInBlacklistError struct {
	ErrCode       string
	LinkExpiresAt time.Time
}

func NewUserSessionsInBlacklistError(code string, linkExpiresAt time.Time) *UserSessionsInBlacklistError {
	return &UserSessionsInBlacklistError{
		ErrCode:       code,
		LinkExpiresAt: linkExpiresAt,
	}
}

func (err *UserSessionsInBlacklistError) Error() string {
	return fmt.Sprintf("sessions in blacklist: error code=%s; link expires at=%s", err.ErrCode, err.LinkExpiresAt)
}

func (err *UserSessionsInBlacklistError) Unwrap() error {
	return ErrUserSessionsInBlackList
}

type LoginFromNewIPOrDeviceError struct {
	ErrCode       string
	CodeExpiresAt time.Time
}

func NewLoginFromNewIPOrDeviceError(code string, expiresAt time.Time) *LoginFromNewIPOrDeviceError {
	return &LoginFromNewIPOrDeviceError{
		ErrCode:       code,
		CodeExpiresAt: expiresAt,
	}
}

func (err *LoginFromNewIPOrDeviceError) Error() string {
	return fmt.Sprintf("login from a new IP address or device: error code=%s; bypass code expires at=%s", err.ErrCode, err.CodeExpiresAt)
}

func (err *LoginFromNewIPOrDeviceError) Unwrap() error {
	return ErrLoginFromNewIPOrDevice
}
