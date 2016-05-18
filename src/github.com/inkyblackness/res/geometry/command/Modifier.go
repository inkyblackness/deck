package command

// Modifier is a unary function to modify a single floating point value.
type Modifier func(float32) float32

// IdentityModifier is the identity function.
func IdentityModifier(value float32) float32 {
	return value
}

// AddingModifier returns a modifier which adds a constant offset to the value.
func AddingModifier(offset float32) Modifier {
	return func(value float32) float32 {
		return value + offset
	}
}
