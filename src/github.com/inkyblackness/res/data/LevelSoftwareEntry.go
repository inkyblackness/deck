package data

import (
	"fmt"
)

// LevelSoftwareEntrySize specifies the byte count of a serialized LevelSoftwareEntry.
const LevelSoftwareEntrySize int = LevelObjectPrefixSize + 3

// LevelSoftwareEntryCount specifies the count how many entries are in one level.
const LevelSoftwareEntryCount int = 16

// LevelSoftwareEntry describes a 'software' level object.
type LevelSoftwareEntry struct {
	LevelObjectPrefix

	ProgramVersion byte
	LogChunkOffset byte
	Unknown        [1]byte
}

// NewLevelSoftwareEntry returns a new LevelSoftwareEntry instance.
func NewLevelSoftwareEntry() *LevelSoftwareEntry {
	return &LevelSoftwareEntry{}
}

func (entry *LevelSoftwareEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()
	result += fmt.Sprintf("Program Version: %d\n", entry.ProgramVersion)
	result += fmt.Sprintf("Log Chunk Offset: %d\n", entry.LogChunkOffset)

	return
}
