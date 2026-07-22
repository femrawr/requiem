package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenString(len int) string {
	bytes := make([]byte, len)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)[:len]
}
