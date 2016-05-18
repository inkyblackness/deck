package image

import (
	"fmt"
	"image/color"
	"io"
)

// SavePalette encodes the provided palette into the given writer.
// Alpha information is lost, the colors are converted to NRGBA and only
// their Red, Green and Blue values are used.
func SavePalette(writer io.Writer, pal color.Palette) (err error) {
	if len(pal) == ColorsPerPixel {
		for _, entry := range pal {
			rgba := color.NRGBAModel.Convert(entry).(color.NRGBA)

			writer.Write([]byte{rgba.R, rgba.G, rgba.B})
		}
	} else {
		err = fmt.Errorf("Palette has wrong length")
	}

	return
}
