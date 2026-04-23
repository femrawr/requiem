package tests_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"shared"
)

const (
	ENCRYPTION_KEY string = "the key to use"

	DATA_TO_ENCRYPT string = "the data to encrypt"

	PAD_ENCRPTION bool = false

	FILE_ENCRYPT_CHUNK_SIZE int = 1 * 1024 * 1024

	FILE_SIZE int = 3.5 * 1024 * 1024
)

func TestCryptData(test *testing.T) {
	encrypted, err := shared.EncryptData(DATA_TO_ENCRYPT, []byte(ENCRYPTION_KEY), PAD_ENCRPTION)
	if err != nil {
		test.Errorf("Failed to encrypt - %v", err)
		return
	}

	fmt.Printf("Encrypted - %s\n", encrypted)

	decrypted, err := shared.DecryptData(encrypted, []byte(ENCRYPTION_KEY), PAD_ENCRPTION)
	if err != nil {
		test.Errorf("Failed to decrypt - %v", err)
		return
	}

	decryptedStr := string(decrypted)
	fmt.Printf("Decrypted - %s\n", decryptedStr)

	if decryptedStr != DATA_TO_ENCRYPT {
		test.Errorf("Decrypted data does not match original data")
		return
	}

	fmt.Println("Decryption successful")
}

func TestCryptFile(test *testing.T) {
	srcFilePath := filepath.Join(os.TempDir(), "test_crypt_file_src")
	defer os.Remove(srcFilePath)

	encFilePath := filepath.Join(os.TempDir(), "crypt_test_enc.bin")
	defer os.Remove(encFilePath)

	decFilePath := filepath.Join(os.TempDir(), "crypt_test_dec.bin")
	defer os.Remove(decFilePath)

	fileData := make([]byte, FILE_SIZE)

	_, err := io.ReadFull(rand.Reader, fileData)
	if err != nil {
		test.Errorf("Failed to generate data - %v", err)
		return
	}

	err = os.WriteFile(srcFilePath, fileData, 0666)
	if err != nil {
		test.Errorf("Failed to write data - %v", err)
		return
	}

	srcFileHash := sha256.Sum256(fileData)
	fmt.Printf("Original file hash  - %x\n", srcFileHash)

	err = shared.EncryptFile(srcFilePath, encFilePath, ENCRYPTION_KEY, FILE_ENCRYPT_CHUNK_SIZE)
	if err != nil {
		test.Errorf("Failed to encrypt file - %v", err)
		return
	}

	err = shared.DecryptFile(encFilePath, decFilePath, ENCRYPTION_KEY, FILE_ENCRYPT_CHUNK_SIZE)
	if err != nil {
		test.Errorf("Failed to decrypt file - %v", err)
		return
	}

	decryptedData, err := os.ReadFile(decFilePath)
	if err != nil {
		test.Errorf("Failed to read data - %v", err)
		return
	}

	decFileHash := sha256.Sum256(decryptedData)
	fmt.Printf("Decrypted hash - %x\n", decFileHash)

	if srcFileHash != decFileHash {
		test.Errorf("Decrypted file hash does not match original file hash")
		return
	}

	fmt.Println("Decryption successful")
}
