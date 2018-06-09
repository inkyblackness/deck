package chunk

import "fmt"

// Identifier represents an integer key of chunks.
type Identifier interface {
	// Value returns the numerical value of the identifier.
	Value() uint16
}

// ID packs an integer as a chunk identifier.
func ID(value uint16) Identifier {
	return chunkID(value)
}

// chunkID identifies a chunk in a resource file.
type chunkID uint16

// Value returns the numerical value of the identifier.
func (id chunkID) Value() uint16 {
	return uint16(id)
}

func (id chunkID) String() string {
	return fmt.Sprintf("%04X", uint16(id))
}
