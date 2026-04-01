package tests_test

import (
	"fmt"
	"testing"

	"requiem/funcs"
)

const VOLUME_TO_SET_TO int = 50

func TestGetMutedDevice(test *testing.T) {
	muted, err := funcs.GetMuted()
	if err != nil {
		test.Errorf("Failed to get muted - %s", err)
		return
	}

	fmt.Printf("Muted - %t\n", muted)
}

func TestSetDeviceVolume(test *testing.T) {
	err := funcs.SetVolume(float32(VOLUME_TO_SET_TO) / 100.0)
	if err != nil {
		test.Errorf("Failed to set volume - %s", err)
		return
	}

	fmt.Printf("Set volume to - %d\n", VOLUME_TO_SET_TO)
}
