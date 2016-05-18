package data

// LevelMarkerEntrySize specifies the byte count of a serialized LevelMarkerEntry.
const LevelMarkerEntrySize int = LevelObjectPrefixSize + 22

// LevelMarkerEntryCount specifies the count how many entries are in one level.
const LevelMarkerEntryCount int = 160

// LevelMarkerEntry describes a 'Marker' level object.
type LevelMarkerEntry struct {
	LevelObjectPrefix

	Unknown [22]byte
}

// NewLevelMarkerEntry returns a new instance of a LevelMarkerEntry.
func NewLevelMarkerEntry() *LevelMarkerEntry {
	return &LevelMarkerEntry{}
}

func (entry *LevelMarkerEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
