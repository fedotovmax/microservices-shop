package domain

import "time"

type DeviceTrustToken struct {
	TokenHash string
	UID       string
	LastUsed  time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (s *DeviceTrustToken) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *DeviceTrustToken) IsRevoked() bool {
	return s.RevokedAt != nil
}

type BlackList struct {
	Code          string
	CodeExpiresAt time.Time
}

type Bypass struct {
	Code            string
	BypassExpiresAt time.Time
}

func (bl *BlackList) IsCodeExpired() bool {
	return time.Now().After(bl.CodeExpiresAt)
}

func (bl *BlackList) ComapreCodes(code string) bool {
	return bl.Code == code
}

func (bp *Bypass) IsCodeExpired() bool {
	return time.Now().After(bp.BypassExpiresAt)
}

func (bp *Bypass) ComapreCodes(code string) bool {
	return bp.Code == code
}

type PreparedTrustTokenAction int8

const (
	TrustTokenNone PreparedTrustTokenAction = iota
	TrustTokenCreated
	TrustTokenUpdated
)

type PreparedTrustToken struct {
	UID                     string
	DeviceTrustTokenValue   string
	DeviceTrustTokenHash    string
	DeviceTrustTokenExpTime time.Time
	Action                  PreparedTrustTokenAction
}
