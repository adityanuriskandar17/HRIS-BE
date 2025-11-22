package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// HashPassword turns a plaintext password into a bcrypt hash using a high cost.
func HashPassword(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("password required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword checks whether the plaintext password matches the given bcrypt hash.
func VerifyPassword(hash, plaintext string) error {
	if hash == "" {
		return errors.New("hash required")
	}
	if plaintext == "" {
		return errors.New("password required")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
