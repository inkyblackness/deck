package model

// Color describes a single color value. The channels have a range of 0..255 .
type Color struct {
	Red   int
	Green int
	Blue  int
}

// RGBA implements the image/color interface.
func (clr Color) RGBA() (r, g, b, a uint32) {
	upscale := func(value int) uint32 {
		b := uint32(value & 0xFF)
		return (b << 8) | b
	}

	return upscale(clr.Red), upscale(clr.Green), upscale(clr.Blue), 0xFFFF
}

// Palette is an identifiable list of colors.
type Palette struct {
	// Colors contains the color values of the palette.
	Colors [256]Color
}
