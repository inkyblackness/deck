package geometry

// VertexContainer provides acces to a list of vertices
type VertexContainer interface {
	// VertexCount returns the number of currently known vertices.
	VertexCount() int
	// Vertex returns the instance for the given index.
	Vertex(index int) Vertex
}
