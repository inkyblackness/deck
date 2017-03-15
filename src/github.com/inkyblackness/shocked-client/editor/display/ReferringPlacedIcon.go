package display

import (
	"github.com/inkyblackness/shocked-client/graphics"
)

type referringPlacedIcon struct {
	center  func() (float32, float32)
	texture func() *graphics.BitmapTexture
}

func (icon *referringPlacedIcon) Center() (x, y float32) {
	return icon.center()
}

func (icon *referringPlacedIcon) Texture() *graphics.BitmapTexture {
	return icon.texture()
}

func (icon *referringPlacedIcon) Size() (width, height float32) {
	return iconSize, iconSize
}
