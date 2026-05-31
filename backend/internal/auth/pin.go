package auth

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPIN(pin string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(pin),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPIN(pin string) bool {
	expected := os.Getenv("PIN_HASH")
	if expected == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(expected),
		[]byte(pin),
	)

	return err == nil
}
