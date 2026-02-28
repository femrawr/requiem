package utils

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"requiem/store"
)

var MutexFile *os.File

func CheckMutex() bool {
	path := filepath.Join(os.TempDir(), store.MUTEX_NAME)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false
	}

	MutexFile = file

	err = lockMutex()
	if err != nil {
		file.Close()
		return false
	}

	return true
}

func RemoveMutex() {
	if MutexFile == nil {
		return
	}

	unlockMutex()
	MutexFile.Close()

	os.Remove(MutexFile.Name())
}

func unlockMutex() error {
	overlap := new(syscall.Overlapped)

	unlocked, _, err := store.UnlockFile.Call(
		uintptr(MutexFile.Fd()),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(overlap)),
	)

	if unlocked == 0 {
		return err
	}

	return nil
}

func lockMutex() error {
	overlap := new(syscall.Overlapped)

	locked, _, err := store.LockFile.Call(
		uintptr(MutexFile.Fd()),
		0x00000002|0x00000001,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(overlap)),
	)

	if locked == 0 {
		return err
	}

	return nil
}
