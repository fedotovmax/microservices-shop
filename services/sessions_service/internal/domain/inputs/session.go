package inputs

import (
	"time"

	"github.com/fedotovmax/grpcutils/violations"
)

type RefreshSession struct {
	userAgent    string
	ip           string
	refreshToken string
}

func NewRefreshSession(refreshToken, userAgent, ip string) *RefreshSession {
	return &RefreshSession{
		refreshToken: refreshToken,
		userAgent:    userAgent,
		ip:           ip,
	}
}

func (i *RefreshSession) GetRefreshToken() string {
	return i.refreshToken
}
func (i *RefreshSession) GetIP() string {
	return i.ip
}

func (i *RefreshSession) GetUserAgent() string {
	return i.userAgent
}

func (i *RefreshSession) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	msg, err := validateIP(i.ip, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("IP", locale, msg, err))
	}

	msg, err = validateEmptyString(i.userAgent, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("UserAgent", locale, msg, err))
	}

	msg, err = validateEmptyString(i.refreshToken, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("RefreshToken", locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

type PrepareSession struct {
	uid              string
	userAgent        string
	ip               string
	bypassCode       *string
	deviceTrustToken *string
}

func NewPrepareSession(uid, userAgent, ip string, code, token *string) *PrepareSession {
	return &PrepareSession{
		uid:              uid,
		userAgent:        userAgent,
		ip:               ip,
		bypassCode:       code,
		deviceTrustToken: token,
	}
}

func (i *PrepareSession) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	msg, err := validateIP(i.ip, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("IP", locale, msg, err))
	}

	msg, err = validateUUID(i.uid, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("UID", locale, msg, err))
	}

	msg, err = validateEmptyString(i.userAgent, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("UserAgent", locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (i *PrepareSession) GetUID() string {
	return i.uid
}
func (i *PrepareSession) GetIP() string {
	return i.ip
}

func (i *PrepareSession) GetUserAgent() string {
	return i.userAgent
}

func (i *PrepareSession) GetBypassCode() string {
	if i.bypassCode != nil {
		return *i.bypassCode
	}
	return ""
}

func (i *PrepareSession) GetDeviceTrustToken() string {
	if i.deviceTrustToken != nil {
		return *i.deviceTrustToken
	}
	return ""
}

type CreateSession struct {
	SID string

	UID string

	RefreshHash string

	Browser        string
	BrowserVersion string

	OS string

	Device string

	IP        string
	ExpiresAt time.Time
}
