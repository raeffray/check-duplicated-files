package hashcalculator

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculatetHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
