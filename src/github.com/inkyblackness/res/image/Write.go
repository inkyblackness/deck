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
func Write(writer io.Writer, bmp Bitmap, bmpType BitmapType, startOffset int) {
	var header BitmapHeader
	pixelData := writePixel(bmp, bmpType)

	header.Type = bmpType
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
	rawPixel := make([]byte, int(bmp.ImageWidth()*bmp.ImageHeight()))

	for row := 0; row < int(bmp.ImageHeight()); row++ {
		copy(rawPixel[int(bmp.ImageWidth())*row:], bmp.Row(row))
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
