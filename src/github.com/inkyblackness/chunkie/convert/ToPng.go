package convert

import (
	"bytes"
	"image/color"
	"image/png"
	"os"

	"github.com/inkyblackness/res/image"
)

// ToPng extracts a bitmap from given block data and saves it to a file.
// The given palette is used should the bitmap not have a private palette.
func ToPng(fileName string, blockData []byte, palette color.Palette) (result bool) {
	bitmap, _ := image.Read(bytes.NewReader(blockData))

	if bitmap != nil {
		img := image.FromBitmap(bitmap, palette)
		file, _ := os.Create(fileName)

		if file != nil {
			defer file.Close()
			png.Encode(file, img)
			result = true
		}
	}

	return
}
