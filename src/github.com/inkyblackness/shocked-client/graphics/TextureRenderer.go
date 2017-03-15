package graphics

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

// TextureRenderer renders textures.
type TextureRenderer interface {
	// Render takes the portion defined by textureRect out of texture to
	// render it within the given display rectangle.
	// textureRect coordinates are given in fractions of the texture.
	Render(modelMatrix *mgl.Mat4, texture Texture, textureRect Rectangle)
}
