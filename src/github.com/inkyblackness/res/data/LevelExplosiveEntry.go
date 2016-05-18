package data

// LevelExplosiveEntrySize specifies the byte count of a serialized LevelExplosiveEntry.
const LevelExplosiveEntrySize int = LevelObjectPrefixSize + 6

// LevelExplosiveEntryCount specifies the count how many entries are in one level.
const LevelExplosiveEntryCount int = 32

// LevelExplosiveEntry describes a 'Explosive' level object.
type LevelExplosiveEntry struct {
	LevelObjectPrefix

	Unknown [6]byte
}

// NewLevelExplosiveEntry returns a new instance of a LevelExplosiveEntry.
func NewLevelExplosiveEntry() *LevelExplosiveEntry {
	return &LevelExplosiveEntry{}
}

func (entry *LevelExplosiveEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
