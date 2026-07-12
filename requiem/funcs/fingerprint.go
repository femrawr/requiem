package funcs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"requiem/utils"

	"golang.org/x/sys/windows/registry"
)

func GenFingerprint() (string, string) {
	fails := 0

	crypto, err := GetCryptoGUID()
	if err != nil {
		fails += 1
	}

	gdid, err := GetIdentityGDID()
	if err != nil {
		fails += 1
	}

	disk, err := GetDiskSerialNumber()
	if err != nil {
		fails += 1
	}

	bios, err := GetBIOSSerialNumber()
	if err != nil {
		fails += 1
	}

	cpu, err := GetProcessorID()
	if err != nil {
		fails += 1
	}

	if fails == 5 {
		return "", ""
	}

	hash := sha256.New()
	hash.Write([]byte(crypto))
	hash.Write([]byte(gdid))
	hash.Write([]byte(disk))
	hash.Write([]byte(bios))
	hash.Write([]byte(cpu))

	theHash := hash.Sum(nil)
	theHMAC := hmac.New(sha256.New, theHash).Sum(nil)

	return hex.EncodeToString(theHash[:]), hex.EncodeToString(theHMAC[:8])
}

func GetCryptoGUID() (string, error) {
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Cryptography",
		registry.QUERY_VALUE|registry.WOW64_64KEY,
	)

	if err != nil {
		return "", err
	}

	val, _, err := key.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}

	key.Close()

	return val, nil
}

func GetIdentityGDID() (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		"SOFTWARE\\Microsoft\\IdentityCRL\\ExtendedProperties",
		registry.QUERY_VALUE,
	)

	if err != nil {
		return "", err
	}

	val, _, err := key.GetStringValue("LID")
	if err != nil {
		return "", err
	}

	key.Close()

	return val, nil
}

func GetDiskSerialNumber() (string, error) {
	return utils.GetCommandOutput("powershell", "(Get-WmiObject Win32_DiskDrive).SerialNumber")
}

func GetBIOSSerialNumber() (string, error) {
	return utils.GetCommandOutput("powershell", "(Get-WmiObject Win32_BIOS).SerialNumber")
}

func GetProcessorID() (string, error) {
	return utils.GetCommandOutput("powershell", "(Get-WmiObject Win32_Processor).ProcessorId")
}
