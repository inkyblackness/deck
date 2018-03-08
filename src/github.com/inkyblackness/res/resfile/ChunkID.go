package resfile

// ChunkID identifies a chunk in a resource file.
type ChunkID uint16

// Value returns the numerical value of the identifier.
func (id ChunkID) Value() uint16 {
	return uint16(id)
}
