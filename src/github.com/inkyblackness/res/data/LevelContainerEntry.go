package data

// LevelContainerEntrySize specifies the byte count of a serialized LevelContainerEntry.
const LevelContainerEntrySize int = LevelObjectPrefixSize + 15

// LevelContainerEntryCount specifies the count how many entries are in one level.
const LevelContainerEntryCount int = 64

// LevelContainerEntry describes a 'Container' level object.
type LevelContainerEntry struct {
	LevelObjectPrefix

	Unknown [15]byte
}

// NewLevelContainerEntry returns a new instance of a LevelContainerEntry.
func NewLevelContainerEntry() *LevelContainerEntry {
	return &LevelContainerEntry{}
}

func (entry *LevelContainerEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
