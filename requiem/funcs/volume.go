package funcs

import (
	"math"
	"syscall"
	"unsafe"

	"requiem/store"
)

var (
	CLSID_MMDeviceEnumerator syscall.GUID = syscall.GUID{
		Data1: 0xBCDE0395,
		Data2: 0xE52F,
		Data3: 0x467C,
		Data4: [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E},
	}

	IID_IMMDeviceEnumerator syscall.GUID = syscall.GUID{
		Data1: 0xA95664D2,
		Data2: 0x9614,
		Data3: 0x4F35,
		Data4: [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6},
	}

	IID_IAudioEndpointVolume syscall.GUID = syscall.GUID{
		Data1: 0x5CDF2C82,
		Data2: 0x841E,
		Data3: 0x4546,
		Data4: [8]byte{0x97, 0x22, 0x0C, 0xF7, 0x40, 0x78, 0x22, 0x9A},
	}
)

type IMMDeviceEnumeratorVTBL struct {
	QueryInterface                         uintptr
	AddRef                                 uintptr
	Release                                uintptr
	EnumAudioEndpoints                     uintptr
	GetDefaultAudioEndpoint                uintptr
	GetDevice                              uintptr
	RegisterEndpointNotificationCallback   uintptr
	UnregisterEndpointNotificationCallback uintptr
}

type IMMDeviceVTBL struct {
	QueryInterface    uintptr
	AddRef            uintptr
	Release           uintptr
	Activate          uintptr
	OpenPropertyStore uintptr
	GetId             uintptr
	GetState          uintptr
}

type IAudioEndpointVolumeVTBL struct {
	QueryInterface                uintptr
	AddRef                        uintptr
	Release                       uintptr
	RegisterControlChangeNotify   uintptr
	UnregisterControlChangeNotify uintptr
	GetChannelCount               uintptr
	SetMasterVolumeLevel          uintptr
	SetMasterVolumeLevelScalar    uintptr
	GetMasterVolumeLevel          uintptr
	GetMasterVolumeLevelScalar    uintptr
}

type IMMDeviceEnumerator struct {
	Vtbl *IMMDeviceEnumeratorVTBL
}

type IMMDevice struct {
	Vtbl *IMMDeviceVTBL
}

type IAudioEndpointVolume struct {
	Vtbl *IAudioEndpointVolumeVTBL
}

func (v *IMMDeviceEnumerator) release() {
	syscall.SyscallN(v.Vtbl.Release, uintptr(unsafe.Pointer(v)))
}

func (v *IMMDevice) release() {
	syscall.SyscallN(v.Vtbl.Release, uintptr(unsafe.Pointer(v)))
}

func (v *IAudioEndpointVolume) release() {
	syscall.SyscallN(v.Vtbl.Release, uintptr(unsafe.Pointer(v)))
}

func SetVolume(volume float32) bool {
	res, _, _ := store.Initialize.Call(0)
	if res != 0 && res != 0x80010106 {
		return false
	}

	defer store.Uninitialize.Call()

	var enumerator *IMMDeviceEnumerator
	res, _, _ = store.Create.Call(
		uintptr(unsafe.Pointer(&CLSID_MMDeviceEnumerator)),
		0,
		23,
		uintptr(unsafe.Pointer(&IID_IMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enumerator)),
	)

	if res != 0 {
		return false
	}

	defer enumerator.release()

	var device *IMMDevice
	res, _, _ = syscall.SyscallN(
		enumerator.Vtbl.GetDefaultAudioEndpoint,
		uintptr(unsafe.Pointer(enumerator)),
		0,
		0,
		uintptr(unsafe.Pointer(&device)),
	)

	if res != 0 {
		return false
	}

	defer device.release()

	var endpoint *IAudioEndpointVolume
	res, _, _ = syscall.SyscallN(
		device.Vtbl.Activate,
		uintptr(unsafe.Pointer(device)),
		uintptr(unsafe.Pointer(&IID_IAudioEndpointVolume)),
		23,
		0,
		uintptr(unsafe.Pointer(&endpoint)),
	)

	if res != 0 {
		return false
	}

	defer endpoint.release()

	res, _, _ = syscall.SyscallN(
		endpoint.Vtbl.SetMasterVolumeLevelScalar,
		uintptr(unsafe.Pointer(endpoint)),
		uintptr(math.Float32bits(volume)),
		0,
	)

	if res != 0 {
		return false
	}

	return true
}
