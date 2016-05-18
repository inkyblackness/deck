package data

// LevelHardwareEntrySize specifies the byte count of a serialized LevelHardwareEntry.
const LevelHardwareEntrySize int = LevelObjectPrefixSize + 1

// LevelHardwareEntryCount specifies the count how many entries are in one level.
const LevelHardwareEntryCount int = 8

// LevelHardwareEntry describes a 'Hardware' level object.
type LevelHardwareEntry struct {
	LevelObjectPrefix

	Unknown [1]byte
}

// NewLevelHardwareEntry returns a new instance of a LevelHardwareEntry.
func NewLevelHardwareEntry() *LevelHardwareEntry {
	return &LevelHardwareEntry{}
}

func (entry *LevelHardwareEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
