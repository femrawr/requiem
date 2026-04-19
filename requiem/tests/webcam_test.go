package tests_test

import (
	"fmt"
	"os"
	"testing"

	"requiem/funcs"
)

const (
	LET_CAMERA_HYDRATE bool   = false
	PICTURE_LOCATION   string = ""
)

func TestTakeWebcamPicture(test *testing.T) {
	pic, err := funcs.TakeWebcam(LET_CAMERA_HYDRATE)
	if err != nil {
		test.Errorf("Failed to capture - %s", err)
		return
	}

	err = os.WriteFile(PICTURE_LOCATION, pic.Bytes(), 0666)
	if err != nil {
		test.Errorf("Failed to write picture - %s", err)
		return
	}

	fmt.Printf("Captured webcam to - %s\n", PICTURE_LOCATION)
}
