package generator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	shareTokenLength = 16
)

func GenerateShareToken() (string, error) {
	bytes := make([]byte, shareTokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
