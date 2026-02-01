package inputs

import "time"

type CreateTrustTokenInput struct {
	TokenHash string
	UID       string
	ExpiresAt time.Time
}

type SecurityInput struct {
	UID           string
	Code          string
	CodeExpiresAt time.Time
}
