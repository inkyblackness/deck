package data

// LevelPatchEntrySize specifies the byte count of a serialized LevelPatchEntry.
const LevelPatchEntrySize int = LevelObjectPrefixSize + 0

// LevelPatchEntryCount specifies the count how many entries are in one level.
const LevelPatchEntryCount int = 32

// LevelPatchEntry describes a 'Patch' level object.
type LevelPatchEntry struct {
	LevelObjectPrefix
}

// NewLevelPatchEntry returns a new instance of a LevelPatchEntry.
func NewLevelPatchEntry() *LevelPatchEntry {
	return &LevelPatchEntry{}
}

func (entry *LevelPatchEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
