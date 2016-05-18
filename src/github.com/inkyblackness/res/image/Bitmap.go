package image

import (
	"image"
	"image/color"
)

// Bitmap wraps the access to a bitmap into abstract functions.
type Bitmap interface {
	// ImageWidth returns the width of the bitmap in pixel.
	ImageWidth() uint16
	// ImageHeight returns the height of the bitmap in pixel.
	ImageHeight() uint16

	// Row returns a slice of the pixel data for given row index.
	Row(index int) []byte

	// Palette returns the private palette of the bitmap. If none is set, this method returns nil.
	Palette() color.Palette

	// Hotspot returns a rectangle within the image bounds. May be 0,0,0,0.
	Hotspot() image.Rectangle
}
