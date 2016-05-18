package data

// LevelSceneryEntrySize specifies the byte count of a serialized LevelSceneryEntry.
const LevelSceneryEntrySize int = LevelObjectPrefixSize + 10

// LevelSceneryEntryCount specifies the count how many entries are in one level.
const LevelSceneryEntryCount int = 176

// LevelSceneryEntry describes an 'Scenery' level object.
type LevelSceneryEntry struct {
	LevelObjectPrefix

	Unknown [10]byte
}

// NewLevelSceneryEntry returns a new LevelSceneryEntry instance.
func NewLevelSceneryEntry() *LevelSceneryEntry {
	return &LevelSceneryEntry{}
}

func (entry *LevelSceneryEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
