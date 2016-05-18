package image

import (
	"image"
	"image/color"
)

// ToBitmap creates a MemoryBitmap out of a paletted image.
// If a palette is provided, this will be used as the private palette of the
// resulting bitmap.
func ToBitmap(img image.PalettedImage, palette color.Palette) *MemoryBitmap {
	var header BitmapHeader

	header.Type = UncompressedBitmap
	header.Width = uint16(img.Bounds().Dx())
	header.Height = uint16(img.Bounds().Dy())
	header.Stride = header.Width
	header.HeightFactor = highestBitShift(header.Height)
	header.WidthFactor = highestBitShift(header.Width)

	data := make([]byte, int(header.Width*header.Height))
	dataOffset := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			data[dataOffset] = img.ColorIndexAt(x, y)
			dataOffset++
		}
	}

	return NewMemoryBitmap(&header, data, palette)
}

func highestBitShift(value uint16) (result byte) {
	if value != 0 {
		for (value >> result) != 1 {
			result++
		}
	}

	return
}
