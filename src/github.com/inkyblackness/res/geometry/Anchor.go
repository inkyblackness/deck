package geometry

// Anchor describes a reference for further nodes or faces
type Anchor interface {
	// Specialize calls the anchor specific handler function on the provided walker.
	Specialize(walker AnchorWalker)
	// Normal returns the normal vector of the anchor.
	Normal() Vector
	// Reference returns the position of the anchor.
	Reference() Vector
}
