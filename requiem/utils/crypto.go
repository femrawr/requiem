package utils

import (
	"requiem/store"
	"shared"
)

func Decrypt(data string) string {
	if data == "" {
		return ""
	}

	dec, err := shared.Decrypt(data, store.CRYPTO_KEY)
	if err != nil {
		return ""
	}

	return dec
}
