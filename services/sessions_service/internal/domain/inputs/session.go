package inputs

import (
	"time"

	"github.com/fedotovmax/grpcutils/violations"
)

type RefreshSessionInput struct {
	userAgent    string
	ip           string
	refreshToken string
}

func NewRefreshSessionInput(refreshToken, userAgent, ip string) *RefreshSessionInput {
	return &RefreshSessionInput{
		refreshToken: refreshToken,
		userAgent:    userAgent,
		ip:           ip,
	}
}

func (i *RefreshSessionInput) GetRefreshToken() string {
	return i.refreshToken
}
func (i *RefreshSessionInput) GetIP() string {
	return i.ip
}

func (i *RefreshSessionInput) GetUserAgent() string {
	return i.userAgent
}

func (i *RefreshSessionInput) Validate(locale string) error {
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

type PrepareSessionInput struct {
	uid              string
	userAgent        string
	ip               string
	bypassCode       *string
	deviceTrustToken *string
}

func NewPrepareSessionInput(uid, userAgent, ip string, code, token *string) *PrepareSessionInput {
	return &PrepareSessionInput{
		uid:              uid,
		userAgent:        userAgent,
		ip:               ip,
		bypassCode:       code,
		deviceTrustToken: token,
	}
}

func (i *PrepareSessionInput) Validate(locale string) error {
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

func (i *PrepareSessionInput) GetUID() string {
	return i.uid
}
func (i *PrepareSessionInput) GetIP() string {
	return i.ip
}

func (i *PrepareSessionInput) GetUserAgent() string {
	return i.userAgent
}

func (i *PrepareSessionInput) GetBypassCode() string {
	if i.bypassCode != nil {
		return *i.bypassCode
	}
	return ""
}

func (i *PrepareSessionInput) GetDeviceTrustToken() string {
	if i.deviceTrustToken != nil {
		return *i.deviceTrustToken
	}
	return ""
}

type SecurityInput struct {
	UID           string
	Code          string
	CodeExpiresAt time.Time
}

type CreateSessionInput struct {
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

type CreateTrustTokenInput struct {
	TokenHash string
	UID       string
	ExpiresAt time.Time
}
