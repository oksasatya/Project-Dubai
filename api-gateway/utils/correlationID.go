package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateCorrelationID generates a unique correlation ID
func GenerateCorrelationID() string {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()) // Fallback jika gagal
	}

	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(randomBytes))
}
