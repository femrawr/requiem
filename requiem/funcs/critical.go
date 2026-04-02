package funcs

import (
	"errors"
	"unsafe"

	"requiem/store"
)

func SetCritical(critical bool) (bool, error) {
	if !store.IsAdmin {
		return false, errors.New("administrator privileges are required to do this")
	}

	var old int32

	ret, _, err := store.AdjustPrivilege.Call(
		uintptr(20),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		return false, err
	}

	if critical {
		ret, _, err = store.SetCritical.Call(uintptr(1), 0, 0)
	} else {
		ret, _, err = store.SetCritical.Call(uintptr(0), 0, 0)
	}

	return ret == 0, err // AAAAAOIJAFIUHASFUISA
}
