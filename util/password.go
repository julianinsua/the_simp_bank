package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Returns the BCypt hash from a given password string
func HashPassword(pswd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// Checks if the provided password matches a given hash
func CheckPassword(pswd, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
}
