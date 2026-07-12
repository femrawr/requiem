package tests_test

import (
	"fmt"
	"testing"

	"requiem/funcs"
)

func TestGenerateFingerprintHash(test *testing.T) {
	hash, hmac := funcs.GenFingerprint()
	if hash == "" || hmac == "" {
		test.Errorf("Failed to generate fingerprint")
		return
	}

	fmt.Printf("Fingerprint hash - %s\n", hash)
	fmt.Printf("Fingerprint hash hmac - %s\n", hmac)
}

func TestGetFingerprintData(test *testing.T) {
	crypto, err := funcs.GetCryptoGUID()
	if err != nil {
		test.Errorf("Failed to get crypto guid - %v", err)
	}

	gdid, err := funcs.GetIdentityGDID()
	if err != nil {
		test.Errorf("Failed to get gdid - %v", err)
	}

	disk, err := funcs.GetDiskSerialNumber()
	if err != nil {
		test.Errorf("Failed to get disk serial - %v", err)
	}

	bios, err := funcs.GetBIOSSerialNumber()
	if err != nil {
		test.Errorf("Failed to get bios serial - %v", err)
	}

	cpu, err := funcs.GetProcessorID()
	if err != nil {
		test.Errorf("Failed to get cpu id - %v", err)
	}

	fmt.Printf("Crypto GUID - %s\n", crypto)
	fmt.Printf("Indentity GDID - %s\n", gdid)
	fmt.Printf("Disk serial - %s\n", disk)
	fmt.Printf("BIOS serial - %s\n", bios)
	fmt.Printf("Processor ID - %s\n", cpu)
}
