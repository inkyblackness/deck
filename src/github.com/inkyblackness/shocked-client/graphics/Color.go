package graphics

// Color represents a 4-element color, ordered R, G, B, A; Range: [0..1]
type Color interface {
	// AsVector returns the color as a 4-dimensional vector.
	AsVector() *[4]float32
}

type colorType [4]float32

// RGBA returns a new color instance with the given R, G, B and A values.
func RGBA(r, g, b, a float32) Color {
	col := &colorType{r, g, b, a}
	return col
}

// AsVector returns the color as a 4-dimensional vector.
func (color *colorType) AsVector() *[4]float32 {
	return (*[4]float32)(color)
}
