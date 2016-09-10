package data

import (
	"fmt"
)

// LevelObjectPrefixSize specifies the byte count of a serialized LevelObjectPrefix.
const LevelObjectPrefixSize int = 6

// LevelObjectPrefix contains the data every Level*Object has at the beginning.
type LevelObjectPrefix struct {
	LevelObjectTableIndex uint16
	Previous              uint16
	Next                  uint16
}

func (prefix *LevelObjectPrefix) String() (result string) {
	result += fmt.Sprintf("Level object table index: %d\n", prefix.LevelObjectTableIndex)
	result += fmt.Sprintf("Links: <- %d | %d ->\n", prefix.Previous, prefix.Next)

	return
}

// NextIndex returns the index of the next object.
func (prefix *LevelObjectPrefix) NextIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(prefix.Next)
}

// SetNextIndex sets the index of the next object.
func (prefix *LevelObjectPrefix) SetNextIndex(index LevelObjectChainIndex) {
	prefix.Next = uint16(index)
}

// PreviousIndex returns the index of the previous object.
func (prefix *LevelObjectPrefix) PreviousIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(prefix.Previous)
}

// SetPreviousIndex sets the index of the previous object.
func (prefix *LevelObjectPrefix) SetPreviousIndex(index LevelObjectChainIndex) {
	prefix.Previous = uint16(index)
}

// ReferenceIndex returns the index of a referenced object.
func (prefix *LevelObjectPrefix) ReferenceIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(prefix.LevelObjectTableIndex)
}

// SetReferenceIndex sets the index of a referenced object.
func (prefix *LevelObjectPrefix) SetReferenceIndex(index LevelObjectChainIndex) {
	prefix.LevelObjectTableIndex = uint16(index)
}
