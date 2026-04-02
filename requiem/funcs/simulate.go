package funcs

import (
	"unsafe"

	"requiem/store"
)

type input struct {
	Type    uint32
	Ki      keyInput
	Padding [8]byte
}

type keyInput struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

func PressUnicodeKey(code uint16) error {
	inputs := [2]input{
		{Type: 1, Ki: keyInput{Scan: code, Flags: 0x0004}},
		{Type: 1, Ki: keyInput{Scan: code, Flags: 0x0004 | 0x0002}},
	}

	ret, _, err := store.SendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret != uintptr(len(inputs)) {
		return err
	}

	return nil
}

func PressVirtualKey(code uint16) error {
	inputs := [2]input{
		{Type: 1, Ki: keyInput{Vk: code}},
		{Type: 1, Ki: keyInput{Vk: code, Flags: 0x0002}},
	}

	ret, _, err := store.SendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret != uintptr(len(inputs)) {
		return err
	}

	return nil
}
