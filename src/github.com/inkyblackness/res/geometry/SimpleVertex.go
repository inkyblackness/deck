package geometry

type simpleVertex struct {
	position Vector
}

// NewSimpleVertex returns a vertex instance based on given parameters.
func NewSimpleVertex(position Vector) Vertex {
	return &simpleVertex{position}
}

func (vertex *simpleVertex) Position() Vector {
	return vertex.position
}
