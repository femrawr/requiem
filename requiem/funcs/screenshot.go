package funcs

import (
	"bytes"
	"image"
	"image/png"
	"unsafe"

	"requiem/store"
)

const (
	_DESKTOPHORZRES uintptr = 118
	_DESKTOPVERTRES uintptr = 117

	_BI_RGB         uintptr = 0
	_DIB_RGB_COLORS uintptr = 0

	_SRCCOPY uintptr = 0x00CC0020
)

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type bitmapInfo struct {
	BmiHeader bitmapInfoHeader
	BmiColors uintptr
}

func CaptureScreen() (*bytes.Buffer, error) {
	context, _, err := store.GetContext.Call(0)
	if context == 0 {
		return nil, err
	}

	defer store.ReleaseCotext.Call(0, context)

	width, _, err := store.GetDeviceCaps.Call(context, _DESKTOPHORZRES)
	if width == 0 {
		return nil, err
	}

	height, _, err := store.GetDeviceCaps.Call(context, _DESKTOPVERTRES)
	if height == 0 {
		return nil, err
	}

	rect := image.Rect(0, 0, int(width), int(height))
	x, y := rect.Dx(), rect.Dy()

	newContext, _, err := store.CreateGoodContext.Call(context)
	if newContext == 0 {
		return nil, err
	}

	defer store.DeleteCotext.Call(newContext)

	info := bitmapInfo{}
	info.BmiHeader.BiSize = uint32(unsafe.Sizeof(info.BmiHeader))
	info.BmiHeader.BiWidth = int32(x)
	info.BmiHeader.BiHeight = int32(-y)
	info.BmiHeader.BiPlanes = 1
	info.BmiHeader.BiBitCount = 32
	info.BmiHeader.BiCompression = uint32(_BI_RGB)

	var bits unsafe.Pointer

	bitmap, _, err := store.CreateBitmapBuffer.Call(
		newContext,
		uintptr(unsafe.Pointer(&info)),
		_DIB_RGB_COLORS,
		uintptr(unsafe.Pointer(&bits)),
		0,
		0,
	)

	if bitmap == 0 {
		return nil, err
	}

	defer store.DeleteObject.Call(bitmap)

	prevObj, _, err := store.SelectObject.Call(newContext, bitmap)
	if prevObj == 0 {
		return nil, err
	}

	defer store.SelectObject.Call(newContext, prevObj)

	ret, _, err := store.BitBlockTransfer.Call(
		newContext,
		0, 0,
		uintptr(x), uintptr(y),
		context,
		0, 0,
		_SRCCOPY,
	)

	if ret == 0 {
		return nil, err
	}

	rawData := unsafe.Slice((*byte)(bits), x*y*4)

	pixels := make([]byte, len(rawData))
	for i := 0; i < len(pixels); i += 4 {
		pixels[i+0] = rawData[i+2]
		pixels[i+1] = rawData[i+1]
		pixels[i+2] = rawData[i+0]
		pixels[i+3] = rawData[i+3]
	}

	picture := &image.RGBA{
		Pix:    pixels,
		Stride: 4 * x,
		Rect:   rect,
	}

	buffer := &bytes.Buffer{}

	err = png.Encode(buffer, picture)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
