package image

import (
	"fmt"
	"image/color"
	"io"
)

// ColorsPerPixel is the constant 256, specifying how many colors are available for one pixel.
const ColorsPerPixel = 256
const bytesPerColor = 3

// LoadPalette reads the raw serialized palette from the given reader and returns the extracted palette.
// The first color has an alpha value of 0x00.
func LoadPalette(reader io.Reader) (pal color.Palette, err error) {
	raw := make([]byte, ColorsPerPixel*bytesPerColor)
	readBytes, err := reader.Read(raw)

	if readBytes == len(raw) {
		pal = make([]color.Color, ColorsPerPixel)
		for i := 0; i < ColorsPerPixel; i++ {
			offset := bytesPerColor * i
			alpha := byte(0xFF)

			if i == 0 {
				alpha = 0x00
			}
			pal[i] = color.NRGBA{R: raw[offset+0], G: raw[offset+1], B: raw[offset+2], A: alpha}
		}
	} else if err == nil {
		err = fmt.Errorf("Too few bytes for palette")
	}

	return
}
