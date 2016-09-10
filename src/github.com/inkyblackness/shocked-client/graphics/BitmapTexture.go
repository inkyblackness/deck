package graphics

import (
	"github.com/inkyblackness/shocked-client/opengl"
)

// BitmapTexture contains a bitmap stored as OpenGL texture.
type BitmapTexture struct {
	gl opengl.OpenGl

	width, height float32
	u, v          float32
	handle        uint32
}

// BitmapRetriever is a thunk that retrieves a cached bitmap.
type BitmapRetriever func() *BitmapTexture

func powerOfTwo(value int) int {
	result := 2

	for (result < value) && (result < 0x1000) {
		result *= 2
	}

	return result
}

// NewBitmapTexture downloads the provided raw data to OpenGL and returns a BitmapTexture instance.
func NewBitmapTexture(gl opengl.OpenGl, width, height int, pixelData []byte) *BitmapTexture {
	textureWidth := powerOfTwo(width)
	textureHeight := powerOfTwo(height)
	tex := &BitmapTexture{
		gl:     gl,
		width:  float32(width),
		height: float32(height),
		handle: gl.GenTextures(1)[0]}
	tex.u = tex.width / float32(textureWidth)
	tex.v = tex.height / float32(textureHeight)

	// The texture has to be blown up to use RGBA from the start;
	// OpenGL 3.2 doesn't know ALPHA format, Open GL ES 2.0 (WebGL) doesn't know RED or R8.
	rgbaData := make([]byte, textureWidth*textureHeight*BytesPerRgba)
	for y := 0; y < height; y++ {
		inStart := y * width
		outOffset := y * textureWidth * BytesPerRgba
		for x := 0; x < width; x++ {
			value := pixelData[inStart+x]
			rgbaData[outOffset+0] = value
			rgbaData[outOffset+1] = value
			rgbaData[outOffset+2] = value
			rgbaData[outOffset+3] = value
			outOffset += BytesPerRgba
		}
	}

	gl.BindTexture(opengl.TEXTURE_2D, tex.handle)
	gl.TexImage2D(opengl.TEXTURE_2D, 0, opengl.RGBA, int32(textureWidth), int32(textureHeight),
		0, opengl.RGBA, opengl.UNSIGNED_BYTE, rgbaData)
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

// Size returns the dimensions of the bitmap, in pixels.
func (tex *BitmapTexture) Size() (width, height float32) {
	return tex.width, tex.height
}

// Handle returns the texture handle.
func (tex *BitmapTexture) Handle() uint32 {
	return tex.handle
}

// UV returns the maximum U and V values for the bitmap. The bitmap will be
// stored in a power-of-two texture, which may be larger than the bitmap.
func (tex *BitmapTexture) UV() (u, v float32) {
	return tex.u, tex.v
}
