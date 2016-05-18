package geometry

// Vector describes a point in three dimensions
type Vector interface {
	// X coordinate
	X() float32
	// Y coordinate
	Y() float32
	// Z coordinate
	Z() float32
}
