package funcs

import (
	"math"
	"syscall"
	"unsafe"

	"requiem/store"
)

var (
	mmDeviceEnumerator = syscall.GUID{
		Data1: 0xBCDE0395,
		Data2: 0xE52F,
		Data3: 0x467C,
		Data4: [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E},
	}

	immDeviceEnumerator = syscall.GUID{
		Data1: 0xA95664D2,
		Data2: 0x9614,
		Data3: 0x4F35,
		Data4: [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6},
	}

	audioEndpointVolume = syscall.GUID{
		Data1: 0x5CDF2C82,
		Data2: 0x841E,
		Data3: 0x4546,
		Data4: [8]byte{0x97, 0x22, 0x0C, 0xF7, 0x40, 0x78, 0x22, 0x9A},
	}
)

type deviceEnumerator struct {
	QueryInterface                         uintptr
	AddRef                                 uintptr
	Release                                uintptr
	EnumAudioEndpoints                     uintptr
	GetDefaultAudioEndpoint                uintptr
	GetDevice                              uintptr
	RegisterEndpointNotificationCallback   uintptr
	UnregisterEndpointNotificationCallback uintptr
}

type device struct {
	QueryInterface    uintptr
	AddRef            uintptr
	Release           uintptr
	Activate          uintptr
	OpenPropertyStore uintptr
	GetId             uintptr
	GetState          uintptr
}

type endpointVolume struct {
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
	SetChannelVolumeLevel         uintptr
	SetChannelVolumeLevelScalar   uintptr
	GetChannelVolumeLevel         uintptr
	GetChannelVolumeLevelScalar   uintptr
	SetMute                       uintptr
	GetMute                       uintptr
}

type deviceEnumeratorStruct struct {
	Table *deviceEnumerator
}

type deviceStruct struct {
	Table *device
}

type endpointVolumeStruct struct {
	Table *endpointVolume
}

func (v *deviceEnumeratorStruct) release() {
	syscall.SyscallN(v.Table.Release, uintptr(unsafe.Pointer(v)))
}

func (v *deviceStruct) release() {
	syscall.SyscallN(v.Table.Release, uintptr(unsafe.Pointer(v)))
}

func (v *endpointVolumeStruct) release() {
	syscall.SyscallN(v.Table.Release, uintptr(unsafe.Pointer(v)))
}

func SetMuted(mute bool) error {
	err := checkInit()
	if err != nil {
		return err
	}

	defer store.Uninitialize.Call()

	enumerator, device, endpoint, err := getDeviceEndpoint()
	if err != nil {
		return err
	}

	defer enumerator.release()
	defer device.release()
	defer endpoint.release()

	muted := 0
	if mute {
		muted = 1
	}

	res, _, err := syscall.SyscallN(
		endpoint.Table.SetMute,
		uintptr(unsafe.Pointer(endpoint)),
		uintptr(muted),
		0,
	)

	if res != 0 {
		return err
	}

	return nil
}

func GetMuted() (bool, error) {
	err := checkInit()
	if err != nil {
		return false, err
	}

	defer store.Uninitialize.Call()

	enumerator, device, endpoint, err := getDeviceEndpoint()
	if err != nil {
		return false, err
	}

	defer enumerator.release()
	defer device.release()
	defer endpoint.release()

	var muted int32
	res, _, err := syscall.SyscallN(
		endpoint.Table.GetMute,
		uintptr(unsafe.Pointer(endpoint)),
		uintptr(unsafe.Pointer(&muted)),
	)

	if res != 0 {
		return false, err
	}

	return muted != 0, nil
}

func SetVolume(volume float32) error {
	err := checkInit()
	if err != nil {
		return err
	}

	defer store.Uninitialize.Call()

	enumerator, device, endpoint, err := getDeviceEndpoint()
	if err != nil {
		return err
	}

	defer enumerator.release()
	defer device.release()
	defer endpoint.release()

	res, _, err := syscall.SyscallN(
		endpoint.Table.SetMasterVolumeLevelScalar,
		uintptr(unsafe.Pointer(endpoint)),
		uintptr(math.Float32bits(volume)),
		0,
	)

	if res != 0 {
		return err
	}

	return nil
}

func checkInit() error {
	res, _, err := store.Initialize.Call(0)
	if res != 0 && res != 0x80010106 {
		return err
	}

	return nil
}

func getDeviceEndpoint() (*deviceEnumeratorStruct, *deviceStruct, *endpointVolumeStruct, error) {
	var enumerator *deviceEnumeratorStruct
	res, _, err := store.Create.Call(
		uintptr(unsafe.Pointer(&mmDeviceEnumerator)),
		0,
		23,
		uintptr(unsafe.Pointer(&immDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enumerator)),
	)

	if res != 0 {
		return nil, nil, nil, err
	}

	var device *deviceStruct
	res, _, err = syscall.SyscallN(
		enumerator.Table.GetDefaultAudioEndpoint,
		uintptr(unsafe.Pointer(enumerator)),
		0,
		0,
		uintptr(unsafe.Pointer(&device)),
	)

	if res != 0 {
		enumerator.release()
		return nil, nil, nil, err
	}

	var endpoint *endpointVolumeStruct
	res, _, err = syscall.SyscallN(
		device.Table.Activate,
		uintptr(unsafe.Pointer(device)),
		uintptr(unsafe.Pointer(&audioEndpointVolume)),
		23,
		0,
		uintptr(unsafe.Pointer(&endpoint)),
	)

	if res != 0 {
		device.release()
		enumerator.release()
		return nil, nil, nil, err
	}

	return enumerator, device, endpoint, nil
}
