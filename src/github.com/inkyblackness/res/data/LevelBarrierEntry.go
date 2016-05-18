package data

// LevelBarrierEntrySize specifies the byte count of a serialized LevelBarrierEntry.
const LevelBarrierEntrySize int = LevelObjectPrefixSize + 8

// LevelBarrierEntryCount specifies the count how many entries are in one level.
const LevelBarrierEntryCount int = 64

// LevelBarrierEntry describes a 'Barrier' level object.
type LevelBarrierEntry struct {
	LevelObjectPrefix

	Unknown [8]byte
}

// NewLevelBarrierEntry returns a new instance of a LevelBarrierEntry.
func NewLevelBarrierEntry() *LevelBarrierEntry {
	return &LevelBarrierEntry{}
}

func (entry *LevelBarrierEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
