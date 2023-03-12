package hashutil

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword return the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// CheckPassword check if the provided password is correct or not
func CheckPassword(password, hashedPassword string) error {
	passwordByte, hashedPasswordByte := []byte(password), []byte(hashedPassword)
	return bcrypt.CompareHashAndPassword(hashedPasswordByte, passwordByte)
}
