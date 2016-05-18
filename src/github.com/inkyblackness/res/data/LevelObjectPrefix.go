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
