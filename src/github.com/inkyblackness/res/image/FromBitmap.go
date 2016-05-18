package image

import (
	"image"
	"image/color"
)

// FromBitmap returns a paletted image instance based on the provided bitmap.
// If the bitmap does not contain a private palette, the provided one is used.
func FromBitmap(bitmap Bitmap, palette color.Palette) image.PalettedImage {
	usedPalette := bitmap.Palette()

	if usedPalette == nil {
		usedPalette = palette
	}

	img := image.NewPaletted(image.Rect(0, 0, int(bitmap.ImageWidth()), int(bitmap.ImageHeight())), usedPalette)
	for row := 0; row < int(bitmap.ImageHeight()); row++ {
		copy(img.Pix[row*img.Stride:], bitmap.Row(row))
	}

	return img
}
