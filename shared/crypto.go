package shared

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func Encrypt(data string, key string) (string, error) {
	hash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, cipher.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	sealed := cipher.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(sealed), nil
}

func Decrypt(data string, key string) (string, error) {
	hash := sha256.Sum256([]byte(key))

	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	size := cipher.NonceSize()
	if len(raw) < size {
		return "", errors.New("data is too short")
	}

	nonce, ciphered := raw[:size], raw[size:]

	result, err := cipher.Open(nil, nonce, ciphered, nil)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
