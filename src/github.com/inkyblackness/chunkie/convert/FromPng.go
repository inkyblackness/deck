package convert

import (
	"image"
	"image/png"
	"os"
)

// FromPng reads a PNG file and encodes it as block.
func FromPng(fileName string, privatePalette bool, compressed bool) []byte {
	file, _ := os.Open(fileName)
	defer file.Close()
	img, _ := png.Decode(file)

	return EncodeImage(img.(*image.Paletted), privatePalette, compressed)
}
