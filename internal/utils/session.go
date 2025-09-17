package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generating session token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
