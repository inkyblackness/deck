package image

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/inkyblackness/res/compress/rle"
)

// Read tries to extract a bitmap from the given source
func Read(source io.ReadSeeker) (bmp Bitmap, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}

	var header BitmapHeader
	var data []byte
	var palette color.Palette = nil

	binary.Read(source, binary.LittleEndian, &header)
	data = make([]byte, int(header.Height)*int(header.Stride))
	if header.Type == CompressedBitmap {
		err = rle.Decompress(source, data)
	} else {
		_, err = source.Read(data)
	}

	if (err == nil) && (header.PaletteOffset != 0) {
		paletteFlag := uint32(0)
		binary.Read(source, binary.LittleEndian, &paletteFlag)
		palette, err = LoadPalette(source)
	}

	if err == nil {
		bmp = NewMemoryBitmap(&header, data, palette)
	}

	return
}
