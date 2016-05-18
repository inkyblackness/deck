package data

// LevelCritterEntrySize specifies the byte count of a serialized LevelCritterEntry.
const LevelCritterEntrySize int = LevelObjectPrefixSize + 40

// LevelCritterEntryCount specifies the count how many entries are in one level.
const LevelCritterEntryCount int = 64

// LevelCritterEntry describes a 'critter' level object.
type LevelCritterEntry struct {
	LevelObjectPrefix

	Unknown [40]byte
}

// NewLevelCritterEntry returns a new LevelCritterEntry instance.
func NewLevelCritterEntry() *LevelCritterEntry {
	return &LevelCritterEntry{}
}

func (entry *LevelCritterEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
