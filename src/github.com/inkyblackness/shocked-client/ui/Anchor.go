package ui

// Anchor provides an absolute value based on some reference. It is used
// for placement and sizing in a visual layout.
type Anchor interface {
	// Value returns the current value of the anchor.
	Value() float32

	// RequestValue suggests a new value for the anchor. Depending on the implementation,
	// the provided value may be taken over, a nearest approximation be used, or
	// ignored alltogether.
	RequestValue(newValue float32)
}
