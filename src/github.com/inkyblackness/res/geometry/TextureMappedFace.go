package geometry

// TextureMappedFace is a face with a texture.
type TextureMappedFace interface {
	Face

	// TextureID returns the identifier of the texture.
	TextureID() uint16
	// TextureCoordinates returns the list of coordinates for the vertices.
	TextureCoordinates() []TextureCoordinate
}
