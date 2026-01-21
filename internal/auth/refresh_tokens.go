package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
