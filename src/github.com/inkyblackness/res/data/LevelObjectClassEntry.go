package data

// LevelObjectClassEntry is a basic entry for a class table, which contains
// the prefix header together with an arbitrary data block.
type LevelObjectClassEntry struct {
	LevelObjectPrefix
	data []byte
}

// NewLevelObjectClassEntry returns a new entry with the given amount as data size.
func NewLevelObjectClassEntry(dataSize int) *LevelObjectClassEntry {
	return &LevelObjectClassEntry{data: make([]byte, dataSize)}
}

// Data returns a slice of the contained data.
func (entry *LevelObjectClassEntry) Data() []byte {
	return entry.data
}
