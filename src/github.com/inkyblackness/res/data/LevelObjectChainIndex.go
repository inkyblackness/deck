package data

// LevelObjectChainIndex describes an index into a chain of objects.
type LevelObjectChainIndex uint16

// LevelObjectChainStartIndex is the table index of the start entry.
const LevelObjectChainStartIndex = LevelObjectChainIndex(0)

// IsStart returns true if the index specifies the starting index.
func (index LevelObjectChainIndex) IsStart() bool {
	return index == LevelObjectChainStartIndex
}
