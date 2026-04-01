package funcs

import (
	"errors"

	"requiem/store"
)

func DisableInputs(disable bool) error {
	if !store.IsAdmin {
		return errors.New("administrator privileges are required to do this")
	}

	disabled := 0
	if disable {
		disabled = 1
	}

	ret, _, err := store.BlockInput.Call(uintptr(disabled))
	if ret == 0 {
		return err
	}

	return nil
}
