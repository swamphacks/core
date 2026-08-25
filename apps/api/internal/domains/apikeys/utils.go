package apikeys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Generates secret of length 64
func generateAPIKeySecret() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func HashAPIKeySecret(secret string) string {
    hash := sha256.Sum256([]byte(secret))
    return hex.EncodeToString(hash[:])
}