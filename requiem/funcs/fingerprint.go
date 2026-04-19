package funcs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"requiem/store"

	"golang.org/x/sys/windows/registry"
)

type storagePropertyQuery struct {
	PropertyId           uint32
	QueryType            uint32
	AdditionalParameters [1]byte
}

type storageDeviceDescriptor struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        byte
	CommandQueueing       byte
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               uint32
	RawPropertiesLength   uint32
	RawDeviceProperties   [1]byte
}

func GenFingerprint() (string, error) {
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

	diskID, err := getDiskSerial()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum224([]byte(val + diskID))
	return hex.EncodeToString(hash[:]), nil
}

func getDiskSerial() (string, error) {
	drive := filepath.VolumeName(store.ExecPath)
	drivePath := "\\\\.\\" + drive

	pointer, err := syscall.UTF16PtrFromString(drivePath)
	if err != nil {
		return "", err
	}

	file, err := syscall.CreateFile(
		pointer,
		0,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	)

	if err != nil {
		return "", err
	}

	defer syscall.CloseHandle(file)

	query := storagePropertyQuery{
		PropertyId: 0,
		QueryType:  0,
	}

	buffer := make([]byte, 512)
	var result uint32

	err = syscall.DeviceIoControl(
		file,
		0x2D1400,
		(*byte)(unsafe.Pointer(&query)),
		uint32(unsafe.Sizeof(query)),
		&buffer[0],
		uint32(len(buffer)),
		&result,
		nil,
	)

	if err != nil {
		return "", err
	}

	if result < uint32(unsafe.Sizeof(storageDeviceDescriptor{})) {
		return "", errors.New("the data returned is too small")
	}

	data := (*storageDeviceDescriptor)(unsafe.Pointer(&buffer[0]))
	if data.SerialNumberOffset == 0 || data.SerialNumberOffset >= uint32(len(buffer)) {
		return "", errors.New("failed to get serial number")
	}

	serial := readCString(buffer[data.SerialNumberOffset:])
	if serial == "" {
		return "", errors.New("serial number empty")
	}

	return strings.TrimSpace(serial), nil
}

func readCString(str []byte) string {
	for i, char := range str {
		if char == 0 {
			return string(str[:i])
		}
	}

	return string(str)
}
