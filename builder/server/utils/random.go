package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenString(len int) string {
	bytes := make([]byte, len)
	rand.Read(bytes)

	return base64.StdEncoding.EncodeToString(bytes)[:len]
}
