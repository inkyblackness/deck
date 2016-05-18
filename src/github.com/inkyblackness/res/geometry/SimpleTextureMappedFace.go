package geometry

type simpleTextureMappedFace struct {
	abstractFace

	textureId          uint16
	textureCoordinates []TextureCoordinate
}

// NewSimpleTextureMappedFace returns a new TextureMappedFace instance with given parameters.
func NewSimpleTextureMappedFace(vertices []int, textureId uint16,
	textureCoordinates []TextureCoordinate) TextureMappedFace {
	return &simpleTextureMappedFace{
		abstractFace:       abstractFace{vertices: vertices},
		textureId:          textureId,
		textureCoordinates: textureCoordinates}
}

func (face *simpleTextureMappedFace) Specialize(walker FaceWalker) {
	walker.TextureMapped(face)
}

func (face *simpleTextureMappedFace) TextureID() uint16 {
	return face.textureId
}

func (face *simpleTextureMappedFace) TextureCoordinates() []TextureCoordinate {
	return face.textureCoordinates
}
