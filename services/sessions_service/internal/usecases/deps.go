package usecases

import "time"

type TokenConfig struct {
	TokenIssuer string
	TokenSecret string

	RefreshExpiresDuration time.Duration
	AccessExpiresDuration  time.Duration
}

type SecurityConfig struct {
	BlacklistCodeExpDuration time.Duration

	LoginBypassExpDuration time.Duration

	DeviceTrustTokenExpDuration time.Duration
	DeviceTrustTokenThreshold   time.Duration

	BlacklistCodeLength uint8

	LoginBypassCodeLength uint8
}
