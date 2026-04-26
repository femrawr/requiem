package tests_test

import (
	"fmt"
	"os"
	"testing"

	"requiem/funcs"
)

const PICTURE_LOCATION string = ""

func TestCaptureWebcam(test *testing.T) {
	pic, err := funcs.CaptureWebcam()
	if err != nil {
		test.Errorf("Failed to capture - %v", err)
		return
	}

	err = os.WriteFile(PICTURE_LOCATION, pic.Bytes(), 0666)
	if err != nil {
		test.Errorf("Failed to write picture - %v", err)
		return
	}

	fmt.Printf("Captured webcam to - %s\n", PICTURE_LOCATION)
}

func TestCaptureScreenshot(test *testing.T) {
	pic, err := funcs.CaptureScreen()
	if err != nil {
		test.Errorf("Failed to capture - %v", err)
		return
	}

	err = os.WriteFile(PICTURE_LOCATION, pic.Bytes(), 0666)
	if err != nil {
		test.Errorf("Failed to write picture - %v", err)
		return
	}

	fmt.Printf("Captured screenshot to - %s\n", PICTURE_LOCATION)
}
