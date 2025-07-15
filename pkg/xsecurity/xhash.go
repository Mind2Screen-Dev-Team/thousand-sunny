package xsecurity

import (
	"crypto/sha256"
	"encoding/hex"
)

func HexHashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
