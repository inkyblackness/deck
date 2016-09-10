package display

import (
	"github.com/inkyblackness/shocked-client/graphics"
)

// PlacedIcon is an icon bitmap with a location.
type PlacedIcon interface {
	Center() (x, y float32)
	Icon() *graphics.BitmapTexture
}
