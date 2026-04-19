package funcs

import (
	"bytes"
	"image/png"

	"github.com/vova616/screenshot"
)

func TakeScreenshot() (*bytes.Buffer, error) {
	ss, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	err = png.Encode(&buffer, ss)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
