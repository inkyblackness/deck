package geometry

type simpleTextureCoordinate struct {
	vertex int
	u      float32
	v      float32
}

// NewSimpleTextureCoordinate returns a TextureCoordinate instance with given parameters.
func NewSimpleTextureCoordinate(vertex int, u, v float32) TextureCoordinate {
	return &simpleTextureCoordinate{
		vertex: vertex,
		u:      u,
		v:      v}
}

func (coord *simpleTextureCoordinate) Vertex() int {
	return coord.vertex
}

func (coord *simpleTextureCoordinate) U() float32 {
	return coord.u
}

func (coord *simpleTextureCoordinate) V() float32 {
	return coord.v
}
