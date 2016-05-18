package graphics

import (
	"github.com/inkyblackness/shocked-client/opengl"
)

// BitmapTexture contains a bitmap stored as OpenGL texture.
type BitmapTexture struct {
	gl opengl.OpenGl

	handle uint32
}

// NewBitmapTexture downloads the provided raw data to OpenGL and returns a BitmapTexture instance.
func NewBitmapTexture(gl opengl.OpenGl, width, height int, pixelData []byte) *BitmapTexture {
	tex := &BitmapTexture{
		gl:     gl,
		handle: gl.GenTextures(1)[0]}

	// The texture has to be blown up to use RGBA from the start;
	// OpenGL 3.2 doesn't know ALPHA format, Open GL ES 1.0 (WebGL) doesn't know RED or R8.
	rgbaData := make([]byte, len(pixelData)*BytesPerRgba)
	for i := 0; i < len(pixelData); i++ {
		value := pixelData[i]
		rgbaData[i*BytesPerRgba+0] = value
		rgbaData[i*BytesPerRgba+1] = value
		rgbaData[i*BytesPerRgba+2] = value
		rgbaData[i*BytesPerRgba+3] = value
	}

	gl.BindTexture(opengl.TEXTURE_2D, tex.handle)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RGBA, int32(width), int32(height), 0, opengl.RGBA, opengl.UNSIGNED_BYTE, rgbaData)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MAG_FILTER, opengl.NEAREST)
	gl.TexParameteri(opengl.TEXTURE_2D, opengl.TEXTURE_MIN_FILTER, opengl.NEAREST)
	gl.GenerateMipmap(opengl.TEXTURE_2D)
	gl.BindTexture(opengl.TEXTURE_2D, 0)

	return tex
}

// Dispose implements the GraphicsTexture interface.
func (tex *BitmapTexture) Dispose() {
	if tex.handle != 0 {
		tex.gl.DeleteTextures([]uint32{tex.handle})
		tex.handle = 0
	}
}

// Handle returns the texture handle.
func (tex *BitmapTexture) Handle() uint32 {
	return tex.handle
}
