package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(plainPassword string) (hashedPassword string, err error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	hashedPassword = string(hashedPasswordBytes)
	return hashedPassword, err
}

// CompareHashAndPassword return nil for success, error for not matched
func CompareHashAndPassword(hashedPassword, plainPassword string) (matched bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return true, err
}
