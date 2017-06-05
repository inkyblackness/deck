package convert

import (
	"bytes"
	goimage "image"

	"github.com/inkyblackness/res/image"
)

// EncodeImage takes a paletted image and encodes it as a block.
func EncodeImage(img *goimage.Paletted, withPrivatePalette bool, compressed, forceTransparency bool) []byte {
	palette := img.Palette
	imgType := image.UncompressedBitmap

	if !withPrivatePalette {
		palette = nil
	}
	if compressed {
		imgType = image.CompressedBitmap
	}
	bmp := image.ToBitmap(img, palette)
	buf := bytes.NewBuffer(nil)
	image.Write(buf, bmp, imgType, forceTransparency, 0)

	return buf.Bytes()
}
