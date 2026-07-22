package shared

import (
	"crypto/sha256"
	"fmt"
	"shared/base"
)

var (
	cryptoKey []byte

	logFunction func(string)
)

func InitKey(key1 string, key2 string) []byte {
	hashed := sha256.Sum256([]byte(key1 + key2))

	cryptoKey = hashed[:]
	return cryptoKey
}

func SetLogFunction(logFunc func(string)) {
	logFunction = logFunc
}

func EncryptConfig(config string) string {
	if config == "" {
		callLogFunction(fmt.Sprintf("failed to encrypt %q - it is empty", config))
		return ""
	}

	if cryptoKey == nil {
		callLogFunction(fmt.Sprintf("failed to encrypt %q - there is no crypto key set", config))
		return ""
	}

	enc, err := base.EncryptData(config, cryptoKey, false)
	if err != nil {
		callLogFunction(fmt.Sprintf("failed to encrypt %q - %v", config, err))
		return ""
	}

	if enc == "" {
		callLogFunction(fmt.Sprintf("encrypted data of %q is empty", config))
	}

	return enc
}

func DecryptConfig(config string) string {
	if config == "" {
		callLogFunction(fmt.Sprintf("failed to decrypt %q - it is empty", config))
		return ""
	}

	if cryptoKey == nil {
		callLogFunction(fmt.Sprintf("failed to decrypt %q - there is no crypto key set", config))
		return ""
	}

	dec, err := base.DecryptData(config, cryptoKey, false)
	if err != nil {
		callLogFunction(fmt.Sprintf("failed to decrypt %q - %v", config, err))
		return ""
	}

	if dec == "" {
		callLogFunction(fmt.Sprintf("decrypted data of %q is empty", config))
	}

	return dec
}

func callLogFunction(data string) {
	if logFunction == nil {
		return
	}

	logFunction(data)
}
