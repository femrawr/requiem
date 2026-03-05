package funcs

import (
	"unsafe"

	"requiem/store"
)

type DeviceModeA struct {
	DmDeviceName         [32]byte
	DmSpecVersion        uint16
	DmDriverVersion      uint16
	DmSize               uint16
	DmDriverExtra        uint16
	DmFields             uint32
	DmPositionX          int32
	DmPositionY          int32
	DmDisplayOrientation uint32
	DmDisplayFixedOutput uint32
	DmColor              int16
	DmDuplex             int16
	DmYResolution        int16
	DmTTOption           int16
	DmCollate            int16
	DmFormName           [32]byte
	DmLogPixels          uint16
	DmBitsPerPel         uint32
	DmPelsWidth          uint32
	DmPelsHeight         uint32
	DmDisplayFlags       uint32
	DmDisplayFrequency   uint32
}

func RotateScreen(by uint32) error {
	var device DeviceModeA
	device.DmSize = uint16(unsafe.Sizeof(device))

	ret, _, err := store.EnumDisplay.Call(
		0,
		uintptr(0xFFFFFFFF),
		uintptr(unsafe.Pointer(&device)),
	)

	if ret == 0 {
		return err
	}

	if (device.DmDisplayOrientation+by)%2 == 1 {
		device.DmPelsWidth, device.DmPelsHeight = device.DmPelsHeight, device.DmPelsWidth
	}

	device.DmDisplayOrientation = by
	device.DmFields = 0x80 | 0x80000 | 0x100000

	ret, _, err = store.ChangeDisplay.Call(
		uintptr(unsafe.Pointer(&device)),
		uintptr(0x01),
	)

	if ret != 0 {
		return err
	}

	return nil
}
