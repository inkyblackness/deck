package geometry

// Vertex describes the properties of a vertex.
type Vertex interface {
	// Position returns the position as a vector.
	Position() Vector
}
