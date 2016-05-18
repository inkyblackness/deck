package geometry

// Face describes a rendered plane.
type Face interface {
	// Specialize calls the face specific handler function on the walker.
	Specialize(walker FaceWalker)
	// Vertices returns the list of indices of the vertices for the face.
	Vertices() []int
}
