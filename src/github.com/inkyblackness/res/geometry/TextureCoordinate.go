package geometry

// TextureCoordinate describes the offset values for a vertex within a texture.
type TextureCoordinate interface {
	// Vertex returns the vertex index of the coordinate
	Vertex() int
	// U offset
	U() float32
	// V offset
	V() float32
}
