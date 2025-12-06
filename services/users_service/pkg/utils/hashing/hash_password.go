package hashing

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Password(password string) (string, error) {

	const op = "hash.Password"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hashed), nil
}

func ComparePassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}
