package graphics

import (
	"github.com/inkyblackness/shocked-client/opengl"
)

// ColorsPerPalette defines how many colors are per palette. This value is 256 to cover byte-based bitmaps.
const ColorsPerPalette = 256

// BytesPerRgba defines the byte count for an RGBA color value.
const BytesPerRgba = 4

// ColorProvider is a function to return the RGBA values for a certain palette index.
type ColorProvider func(index int) (byte, byte, byte, byte)

// PaletteTexture contains a palette stored as OpenGL texture.
type PaletteTexture struct {
	gl opengl.OpenGl

	colorProvider ColorProvider
	handle        uint32
}

// NewPaletteTexture creates a new PaletteTexture instance.
func NewPaletteTexture(gl opengl.OpenGl, colorProvider ColorProvider) *PaletteTexture {
	tex := &PaletteTexture{
		gl:            gl,
		colorProvider: colorProvider,
		handle:        gl.GenTextures(1)[0]}

	tex.Update()

	return tex
}

// Dispose implements the GraphicsTexture interface.
func (tex *PaletteTexture) Dispose() {
	if tex.handle != 0 {
		tex.gl.DeleteTextures([]uint32{tex.handle})
		tex.handle = 0
	}
}

// Handle returns the texture handle.
func (tex *PaletteTexture) Handle() uint32 {
	return tex.handle
}

// Update reloads the palette.
func (tex *PaletteTexture) Update() {
	gl := tex.gl
	var palette [ColorsPerPalette * BytesPerRgba]byte

	tex.loadColors(&palette)
	gl.BindTexture(opengl.TEXTURE_2D, tex.handle)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RGBA, ColorsPerPalette, 1, 0, opengl.RGBA, opengl.UNSIGNED_BYTE, palette[:])
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.NEAREST)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.NEAREST)
	gl.GenerateMipmap(opengl.TEXTURE_2D)
	gl.BindTexture(opengl.TEXTURE_2D, 0)
}

func (tex *PaletteTexture) loadColors(palette *[ColorsPerPalette * BytesPerRgba]byte) {
	for i := 0; i < ColorsPerPalette; i++ {
		r, g, b, a := tex.colorProvider(i)

		palette[i*BytesPerRgba+0] = r
		palette[i*BytesPerRgba+1] = g
		palette[i*BytesPerRgba+2] = b
		palette[i*BytesPerRgba+3] = a
	}
}
