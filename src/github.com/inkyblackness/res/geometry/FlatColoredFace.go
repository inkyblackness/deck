package geometry

// FlatColoredFace describes a face in a single color.
type FlatColoredFace interface {
	Face

	// Color returns the color index of the face.
	Color() ColorIndex
}
