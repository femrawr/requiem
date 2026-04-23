package shared

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func EncryptData(data string, key []byte, pad bool) (string, error) {
	keyHash := sha256.Sum256(key)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	sealed := gcm.Seal(nonce, nonce, []byte(data), nil)
	encoded := base64.StdEncoding.EncodeToString(sealed)

	return encoded, nil
}

func DecryptData(data string, key []byte, padded bool) (string, error) {
	keyHash := sha256.Sum256(key)

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	size := gcm.NonceSize()
	if len(decoded) < size {
		return "", errors.New("data is too short")
	}

	result, err := gcm.Open(nil, decoded[:size], decoded[size:], nil)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func EncryptFile(srcFilePath string, outFilePath string, key string, chunkSize int) error {
	src, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}

	defer src.Close()

	out, err := os.Create(outFilePath)
	if err != nil {
		return err
	}

	defer out.Close()

	keyHash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	buffer := make([]byte, chunkSize)

	for {
		length, err := src.Read(buffer)
		if length > 0 {
			nonce := make([]byte, gcm.NonceSize())

			_, err := io.ReadFull(rand.Reader, nonce)
			if err != nil {
				return err
			}

			encrypted := gcm.Seal(nonce, nonce, buffer[:length], nil)

			_, err = out.Write(encrypted)
			if err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func DecryptFile(srcFilePath string, outFilePath string, key string, chunkSize int) error {
	src, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}

	defer src.Close()

	out, err := os.Create(outFilePath)
	if err != nil {
		return err
	}

	defer out.Close()

	keyHash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	encChunkSize := gcm.NonceSize() + chunkSize + gcm.Overhead()
	buffer := make([]byte, encChunkSize)

	for {
		length, err := io.ReadFull(src, buffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if length == 0 {
				break
			}
		} else if err != nil {
			return err
		}

		chunk := buffer[:length]
		if len(chunk) < gcm.NonceSize() {
			return errors.New("chunk too short")
		}

		nonce, ciphered := chunk[:gcm.NonceSize()], chunk[gcm.NonceSize():]

		decrypted, err := gcm.Open(nil, nonce, ciphered, nil)
		if err != nil {
			return err
		}

		_, err = out.Write(decrypted)
		if err != nil {
			return err
		}

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}

	return nil
}
