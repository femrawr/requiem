package utils

import "unsafe"

func ReadStringFromMemory(pointer uintptr) string {
	if pointer == 0 {
		return ""
	}

	var buffer []byte

	for {
		byte := *(*byte)(unsafe.Pointer(pointer))
		if byte == 0 {
			break
		}

		buffer = append(buffer, byte)
		pointer++
	}

	return string(buffer)
}
