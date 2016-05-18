package geometry

// ShadeColoredFace is a face with a color and a specific shade.
type ShadeColoredFace interface {
	Face

	// Color returns the color index of the face.
	Color() ColorIndex
	// Shade returns the shade intensity of the color.
	Shade() uint16
}
