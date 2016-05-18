package data

// LevelItemEntrySize specifies the byte count of a serialized LevelItemEntry.
const LevelItemEntrySize int = LevelObjectPrefixSize + 10

// LevelItemEntryCount specifies the count how many entries are in one level.
const LevelItemEntryCount int = 128

// LevelItemEntry describes an 'item' level object.
type LevelItemEntry struct {
	LevelObjectPrefix

	Unknown [10]byte
}

// NewLevelItemEntry returns a new LevelItemEntry instance.
func NewLevelItemEntry() *LevelItemEntry {
	return &LevelItemEntry{}
}

func (entry *LevelItemEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
