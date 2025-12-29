package inputs

import (
	"time"

	"github.com/fedotovmax/grpcutils/violations"
)

type VerifyAccessInput struct {
	issuer      string
	accessToken string
}

func NewVerifyAccessInput(accessToken, issuer string) *VerifyAccessInput {
	return &VerifyAccessInput{
		accessToken: accessToken,
		issuer:      issuer,
	}
}
func (i *VerifyAccessInput) GetAccessToken() string {
	return i.accessToken
}
func (i *VerifyAccessInput) GetIssuer() string {
	return i.issuer
}

func (i *VerifyAccessInput) Validate(locale string) error {
	var validationErrors violations.ValidationErrors

	msg, err := validateEmptyString(i.issuer, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("Issuer", locale, msg, err))
	}

	msg, err = validateEmptyString(i.accessToken, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("AccessToken", locale, msg, err))
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

type RefreshSessionInput struct {
	userAgent    string
	ip           string
	issuer       string
	refreshToken string
}

func NewRefreshSessionInput(refreshToken, userAgent, ip, issuer string) *RefreshSessionInput {
	return &RefreshSessionInput{
		refreshToken: refreshToken,
		userAgent:    userAgent,
		ip:           ip,
		issuer:       issuer,
	}
}

func (i *RefreshSessionInput) GetRefreshToken() string {
	return i.refreshToken
}
func (i *RefreshSessionInput) GetIP() string {
	return i.ip
}
func (i *RefreshSessionInput) GetIssuer() string {
	return i.issuer
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

	msg, err = validateEmptyString(i.issuer, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("Issuer", locale, msg, err))
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
	uid       string
	userAgent string
	ip        string
	issuer    string
}

func NewPrepareSessionInput(uid, userAgent, ip, issuer string) *PrepareSessionInput {
	return &PrepareSessionInput{
		uid:       uid,
		userAgent: userAgent,
		ip:        ip,
		issuer:    issuer,
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

	msg, err = validateEmptyString(i.issuer, locale)

	if err != nil {
		validationErrors = append(validationErrors, addValidationError("Issuer", locale, msg, err))
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
func (i *PrepareSessionInput) GetIssuer() string {
	return i.issuer
}
func (i *PrepareSessionInput) GetUserAgent() string {
	return i.userAgent
}

type AddToBlackListInput struct {
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
