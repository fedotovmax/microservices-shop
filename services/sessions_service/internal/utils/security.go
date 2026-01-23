package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"time"
)

type CreateTokenResult struct {
	Nohashed string
	Hashed   string
}

func CreateToken() (*CreateTokenResult, error) {
	token, err := GenerateToken()

	if err != nil {
		return nil, err
	}

	hash := HashToken(token)

	resulst := &CreateTokenResult{
		Nohashed: token,
		Hashed:   hash,
	}

	return resulst, nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func GenerateSecurityCode(length uint8) (string, error) {

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}

func ExtendTrustTokenTTL(currentExp time.Time, nowUTC time.Time, threshold time.Duration, maxTTL time.Duration) time.Time {
	remaining := currentExp.Sub(nowUTC)

	if remaining <= threshold {
		return nowUTC.Add(maxTTL)
	}

	return currentExp
}
