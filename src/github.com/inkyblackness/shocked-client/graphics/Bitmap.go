package graphics

import (
	"encoding/base64"

	"github.com/inkyblackness/shocked-model"
)

// Bitmap is a simple palette based image.
type Bitmap struct {
	Width  int
	Height int
	Pixels []byte
}

// BitmapFromRaw returns a bitmap instance with decoded pixel data.
func BitmapFromRaw(raw model.RawBitmap) Bitmap {
	pixelData, _ := base64.StdEncoding.DecodeString(raw.Pixels)

	return Bitmap{
		Width:  raw.Width,
		Height: raw.Height,
		Pixels: pixelData}
}
