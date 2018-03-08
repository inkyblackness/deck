package resfile

// Identifier represents an integral key of chunks.
type Identifier interface {
	// Value returns the numerical value of the identifier.
	Value() uint16
}
