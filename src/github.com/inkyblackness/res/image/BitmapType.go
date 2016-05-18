package image

import "fmt"

type BitmapType int16

const (
	UncompressedBitmap BitmapType = 0x0002
	CompressedBitmap   BitmapType = 0x0004
)

func (bmpType BitmapType) String() (result string) {
	switch bmpType {
	case UncompressedBitmap:
		result = "Uncompressed"
	case CompressedBitmap:
		result = "Compressed"
	default:
		result = fmt.Sprintf("Unknown (0x%04X)", int16(bmpType))
	}

	return
}
