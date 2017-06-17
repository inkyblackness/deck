package graphics

import (
	"image"
	"image/color"
	"math"
)

// reference white point
var d65 = [3]float64{0.95047, 1.00000, 1.08883}

func labF(t float64) float64 {
	if t > 6.0/29.0*6.0/29.0*6.0/29.0 {
		return math.Cbrt(t)
	}
	return t/3.0*29.0/6.0*29.0/6.0 + 4.0/29.0
}

func square(value float64) float64 {
	return value * value
}

type labEntry struct {
	l float64
	a float64
	b float64
}

func labEntryFromColor(clr color.Color) labEntry {
	rLinear, gLinear, bLinear, _ := clr.RGBA()
	r, g, b := float64(rLinear)/float64(0xFFFF), float64(gLinear)/float64(0xFFFF), float64(bLinear)/float64(0xFFFF)
	x := 0.4124564*r + 0.3575761*g + 0.1804375*b
	y := 0.2126729*r + 0.7151522*g + 0.0721750*b
	z := 0.0193339*r + 0.1191920*g + 0.9503041*b
	whiteRef := d65
	fy := labF(y / whiteRef[1])
	entry := labEntry{
		l: 1.16*fy - 0.16,
		a: 5.0 * (labF(x/whiteRef[0]) - fy),
		b: 2.0 * (fy - labF(z/whiteRef[2]))}

	return entry
}

func (entry labEntry) distanceTo(other labEntry) float64 {
	return math.Sqrt(square(entry.l-other.l) + square(entry.a-other.a) + square(entry.b-other.b))
}

// StandardBitmapper creates bitmap images from generic images.
type StandardBitmapper struct {
	pal []labEntry
}

// NewStandardBitmapper returns a new bitmapper instance.
func NewStandardBitmapper(palette []color.Color) *StandardBitmapper {
	bitmapper := &StandardBitmapper{}

	for _, clr := range palette {
		bitmapper.pal = append(bitmapper.pal, labEntryFromColor(clr))
	}

	return bitmapper
}

// Map maps the provided image to a bitmap based on the internal palette.
func (bitmapper *StandardBitmapper) Map(img image.Image) Bitmap {
	var bmp Bitmap
	bounds := img.Bounds()

	bmp.Width = bounds.Dx()
	bmp.Height = bounds.Dy()
	bmp.Pixels = make([]byte, bmp.Width*bmp.Height)
	for row := 0; row < bmp.Height; row++ {
		for column := 0; column < bmp.Width; column++ {
			bmp.Pixels[row*bmp.Width+column] = bitmapper.MapColor(img.At(column, row))
		}
	}

	return bmp
}

// MapColor maps the provided color to the nearest index in the palette.
func (bitmapper *StandardBitmapper) MapColor(clr color.Color) (palIndex byte) {
	_, _, _, a := clr.RGBA()

	if a > 0 {
		clrEntry := labEntryFromColor(clr)
		palDistance := 1000.0

		for colorIndex, palEntry := range bitmapper.pal {
			distance := palEntry.distanceTo(clrEntry)
			if distance < palDistance {
				palDistance = distance
				palIndex = byte(colorIndex)
			}
		}
	}
	return
}
