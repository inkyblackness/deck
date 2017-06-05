package image

import (
	"fmt"
)

// BitmapHeader describes the header of an encoded bitmap
type BitmapHeader struct {
	Unknown0000      [4]byte
	Type             BitmapType
	TransparencyFlag int16

	Width        uint16
	Height       uint16
	Stride       uint16
	WidthFactor  byte
	HeightFactor byte
	HotspotBox   [4]uint16

	PaletteOffset int32
}

func (header *BitmapHeader) String() (result string) {
	result += fmt.Sprintf("Type: %v, %dx%d\n", header.Type, header.Width, header.Height)
	result += fmt.Sprintf("Transparency: %v\n", header.TransparencyFlag != 0)
	result += fmt.Sprintf("%d,%d | %d,%d\n", header.HotspotBox[0], header.HotspotBox[1], header.HotspotBox[2], header.HotspotBox[3])
	result += fmt.Sprintf("PaletteOffset: %d", header.PaletteOffset)
	return
}
