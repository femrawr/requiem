package funcs

import (
	"bytes"
	"image/jpeg"

	"github.com/vova616/screenshot"
)

func TakeScreenshot() *bytes.Buffer {
	ss, err := screenshot.CaptureScreen()
	if err != nil {
		return nil
	}

	var buffer bytes.Buffer

	err = jpeg.Encode(&buffer, ss, nil)
	if err != nil {
		return nil
	}

	return &buffer
}
