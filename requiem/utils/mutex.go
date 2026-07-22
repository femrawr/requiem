package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"requiem/store"
)

const (
	_DEBUG_MUTEX_PREFIX string = "requiem mutex"

	_LOCKFILE_EXCLUSIVE_LOCK   uintptr = 0x00000002
	_LOCKFILE_FAIL_IMMEDIATELY uintptr = 0x00000001
)

var mutexFile *os.File

func CheckMutex() bool {
	mutexName := store.MUTEX_NAME
	if store.DEBUG_MODE {
		mutexName = fmt.Sprintf("%s %s", _DEBUG_MUTEX_PREFIX, mutexName)
	}

	path := filepath.Join(os.TempDir(), mutexName)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false
	}

	mutexFile = file

	err = lockMutex()
	if err != nil {
		file.Close()
		return false
	}

	return true
}

func RemoveMutex() {
	if mutexFile == nil {
		return
	}

	unlockMutex()
	mutexFile.Close()

	os.Remove(mutexFile.Name())
}

func unlockMutex() error {
	overlap := new(syscall.Overlapped)

	unlocked, _, err := store.UnlockFile.Call(
		uintptr(mutexFile.Fd()),
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
		uintptr(mutexFile.Fd()),
		_LOCKFILE_EXCLUSIVE_LOCK|_LOCKFILE_FAIL_IMMEDIATELY,
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
