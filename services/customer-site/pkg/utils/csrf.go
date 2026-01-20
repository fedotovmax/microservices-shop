package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func NewCSRF() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
