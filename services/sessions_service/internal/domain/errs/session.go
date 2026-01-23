package errs

import (
	"errors"
)

var ErrTrustTokenNotFound = errors.New("trust token not found or expired or revoked")

var ErrUserDeleted = errors.New("user is deleted")

var ErrInvalidSession = errors.New("invalid session")

var ErrSessionNotFound = errors.New("session not found")

var ErrSessionExpired = errors.New("session is expired")

var ErrSessionRevoked = errors.New("session is revoked")

var ErrUserSessionsInBlackList = errors.New("session in blacklist")

var ErrBlacklistCodeExpired = errors.New("blacklist code expired")

var ErrBadBlacklistCode = errors.New("bad blacklist code")

var ErrAgentLooksLikeBot = errors.New("the user agent looks like a bot")

var ErrInvalidSessionIP = errors.New("invalid session IP")

var ErrLoginFromNewIPOrDevice = errors.New("login from a new IP address or device")

var ErrBypassCodeExpired = errors.New("bypass code expired")

var ErrBadBypassCode = errors.New("bad bypass code")

// type UserSessionRevokedError struct {
// 	Email string
// 	UID   string
// 	SID   string
// }

// func NewUserSessionRevokedError(email, uid, sid string) *UserSessionRevokedError {
// 	return &UserSessionRevokedError{
// 		Email: email,
// 		UID:   uid,
// 		SID:   sid,
// 	}
// }

// func (err *UserSessionRevokedError) Error() string {
// 	return fmt.Sprintf("session revoked: sid=%s; uid=%s; email=%s", err.SID, err.UID, err.Email)
// }

// func (err *UserSessionRevokedError) Unwrap() error {
// 	return ErrSessionRevoked
// }

// type UserSessionsInBlacklistError struct {
// 	Email              string
// 	UID                string
// 	NeedNewUnblockCode bool
// }

// func NewUserSessionsInBlacklistError(email, uid string) *UserSessionsInBlacklistError {
// 	return &UserSessionsInBlacklistError{
// 		Email: email,
// 		UID:   uid,
// 	}
// }

// func (err *UserSessionsInBlacklistError) Error() string {
// 	return fmt.Sprintf("sessions in blacklist: uid=%s; email=%s", err.UID, err.Email)
// }

// func (err *UserSessionsInBlacklistError) Unwrap() error {
// 	return ErrUserSessionsInBlackList
// }

// type LoginFromNewIPOrDeviceError struct {
// 	Email string
// 	UID   string
// }

// func NewLoginFromNewIPOrDeviceError(email, uid string) *LoginFromNewIPOrDeviceError {
// 	return &LoginFromNewIPOrDeviceError{
// 		Email: email,
// 		UID:   uid,
// 	}
// }

// func (err *LoginFromNewIPOrDeviceError) Error() string {
// 	return fmt.Sprintf("login from a new IP address or device: uid=%s; email=%s", err.UID, err.Email)
// }

// func (err *LoginFromNewIPOrDeviceError) Unwrap() error {
// 	return ErrLoginFromNewIPOrDevice
// }
