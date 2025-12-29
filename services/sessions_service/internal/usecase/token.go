package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type createRefreshTokenResult struct {
	nohashed string
	hashed   string
}

func (u *usecases) createRefreshToken() (*createRefreshTokenResult, error) {
	refreshToken, err := u.generateRefreshToken()

	if err != nil {
		return nil, err
	}

	refreshHash := u.hashToken(refreshToken)

	resulst := &createRefreshTokenResult{
		nohashed: refreshToken,
		hashed:   refreshHash,
	}

	return resulst, nil
}

func (u *usecases) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (u *usecases) hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// func (u *usecases) compareHashes(hashed, noHashed string) bool {
// 	clientHash := u.hashToken(noHashed)
// 	return hmac.Equal([]byte(clientHash), []byte(hashed))
// }
