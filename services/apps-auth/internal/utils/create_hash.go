package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func CreateHash(str string) string {
	sum := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}
