package funcs

import (
	"errors"
	"requiem/store"
)

func DisableInputs(disable bool) (bool, error) {
	if !store.IsAdmin {
		return false, errors.New("administrator privileges are required to do this")
	}

	block := 0
	if disable {
		block = 1
	}

	ret, _, err := store.BlockInput.Call(uintptr(block))
	if ret == 0 {
		return false, err
	}

	return true, nil
}
