package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// HashPIN returns a SHA-256 hex digest of the PIN.
// Use this once to generate the value to put in your .env.
func HashPIN(pin string) string {
	h := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(h[:])
}

// VerifyPIN checks the supplied PIN against the hash in PIN_HASH env var.
func VerifyPIN(pin string) bool {
	expected := os.Getenv("PIN_HASH")
	if expected == "" {
		return false
	}
	return HashPIN(pin) == expected
}
