package inputs

import "time"

type CreateTrustToken struct {
	TokenHash string
	UID       string
	ExpiresAt time.Time
}

type Security struct {
	UID           string
	Code          string
	CodeExpiresAt time.Time
}
