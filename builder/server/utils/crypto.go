package utils

import (
	"fmt"
	"shared"
)

var cryptoKey string

func SetCryptKey(key string) {
	cryptoKey = key
}

func Encrypt(data string) string {
	if data == "" {
		return ""
	}

	enc, err := shared.Encrypt(data, cryptoKey)
	if err != nil {
		fmt.Printf("failed to encrypt %s - %s\n", data, err)
		return ""
	}

	return enc
}
