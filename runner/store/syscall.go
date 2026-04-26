package store

import (
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ntdll    = syscall.NewLazyDLL("ntdll.dll")
)

var (
	getProcAddress = kernel32.NewProc("GetProcAddress")

	VirtualAlloc   = kernel32.NewProc("VirtualAlloc")
	VirtualProtect = kernel32.NewProc("VirtualProtect")
	LoadLibrary    = kernel32.NewProc("LoadLibraryA")

	AddFunctionTable = ntdll.NewProc("RtlAddFunctionTable")
)

func GetFunctionAddress(module uintptr, name string) (uintptr, error) {
	pointer, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, err
	}

	address, _, err := getProcAddress.Call(module, uintptr(unsafe.Pointer(pointer)))
	if address == 0 {
		return 0, err
	}

	return address, nil
}
