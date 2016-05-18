package data

// LevelAnimationEntrySize specifies the byte count of a serialized LevelAnimationEntry.
const LevelAnimationEntrySize int = LevelObjectPrefixSize + 4

// LevelAnimationEntryCount specifies the count how many entries are in one level.
const LevelAnimationEntryCount int = 32

// LevelAnimationEntry describes a 'Animation' level object.
type LevelAnimationEntry struct {
	LevelObjectPrefix

	Unknown [4]byte
}

// NewLevelAnimationEntry returns a new instance of a LevelAnimationEntry.
func NewLevelAnimationEntry() *LevelAnimationEntry {
	return &LevelAnimationEntry{}
}

func (entry *LevelAnimationEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
