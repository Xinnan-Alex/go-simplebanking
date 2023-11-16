package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	haspassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(haspassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
