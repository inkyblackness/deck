package image

import (
	"image"
	"image/color"
)

type MemoryBitmap struct {
	header  BitmapHeader
	data    []byte
	palette color.Palette
}

// NewMemoryBitmap returns a MemoryBitmap instance based on the provided data.
func NewMemoryBitmap(header *BitmapHeader, data []byte, palette color.Palette) *MemoryBitmap {
	bmp := &MemoryBitmap{
		header:  *header,
		data:    data,
		palette: palette}

	return bmp
}

// Compressed returns whether the bitmap data is stored in compressed form
func (bmp *MemoryBitmap) Compressed() bool {
	return bmp.header.Type == CompressedBitmap
}

// ImageWidth returns the width of the bitmap in pixel.
func (bmp *MemoryBitmap) ImageWidth() uint16 {
	return bmp.header.Width
}

// ImageHeight returns the height of the bitmap in pixel.
func (bmp *MemoryBitmap) ImageHeight() uint16 {
	return bmp.header.Height
}

// Row returns a slice of the pixel data for given row index.
func (bmp *MemoryBitmap) Row(index int) []byte {
	start := index * int(bmp.header.Stride)

	return bmp.data[start : start+int(bmp.header.Stride)]
}

// Palette returns the private palette of the bitmap. If none is set, this method returns nil.
func (bmp *MemoryBitmap) Palette() color.Palette {
	return bmp.palette
}

// Hotspot returns a rectangle within the image bounds. May be 0,0,0,0.
func (bmp *MemoryBitmap) Hotspot() image.Rectangle {
	return image.Rect(int(bmp.header.HotspotBox[0]), int(bmp.header.HotspotBox[1]),
		int(bmp.header.HotspotBox[2]), int(bmp.header.HotspotBox[3]))
}
