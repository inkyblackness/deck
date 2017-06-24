package image

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/inkyblackness/res/compress/rle"
)

// Write serializes given bitmap. The provided type specifies whether the pixel data
// shall be compressed. The start offset is used for images with a private palette.
// The start offset specifies the byte length of the blocks before the current one.
// Force transparency indicates that palette index 0x00 is meant to be treated as
// transparent. Some usages of bitmaps imply transparency and don't require
// this flag to be set.
func Write(writer io.Writer, bmp Bitmap, bmpType BitmapType, forceTransparency bool, startOffset int) {
	var header BitmapHeader
	pixelData := writePixel(bmp, bmpType)

	header.Type = bmpType
	if forceTransparency {
		header.TransparencyFlag = 1
	}
	header.Width = bmp.ImageWidth()
	header.Height = bmp.ImageHeight()
	header.Stride = header.Width
	header.HeightFactor = highestBitShift(header.Height)
	header.WidthFactor = highestBitShift(header.Width)
	if bmp.Palette() != nil {
		header.PaletteOffset = int32(startOffset + len(pixelData) + binary.Size(header))
	}

	binary.Write(writer, binary.LittleEndian, &header)
	writer.Write(pixelData)
	if bmp.Palette() != nil {
		binary.Write(writer, binary.LittleEndian, privatePaletteFlag)
		SavePalette(writer, bmp.Palette())
	}
}

func writePixel(bmp Bitmap, bmpType BitmapType) (result []byte) {
	width := int(bmp.ImageWidth())
	height := int(bmp.ImageHeight())
	rawPixel := make([]byte, width*height)

	for row := 0; row < height; row++ {
		inRow := bmp.Row(row)
		outRow := rawPixel[width*row:]
		copy(outRow, inRow)
	}

	if bmpType == CompressedBitmap {
		buf := bytes.NewBuffer(nil)
		rle.Compress(buf, rawPixel)
		result = buf.Bytes()
	} else {
		result = rawPixel
	}

	return
}
