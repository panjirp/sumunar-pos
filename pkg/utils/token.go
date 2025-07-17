package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken creates a secure 256-bit token (32 bytes).
func GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// fallback: this should rarely happen
		return ""
	}
	return hex.EncodeToString(b)
}
