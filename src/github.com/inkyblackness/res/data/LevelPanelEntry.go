package data

// LevelPanelEntrySize specifies the byte count of a serialized LevelPanelEntry.
const LevelPanelEntrySize int = LevelObjectPrefixSize + 24

// LevelPanelEntryCount specifies the count how many entries are in one level.
const LevelPanelEntryCount int = 64

// LevelPanelEntry describes a 'Panel' level object.
type LevelPanelEntry struct {
	LevelObjectPrefix

	Unknown [24]byte
}

// NewLevelPanelEntry returns a new instance of a LevelPanelEntry.
func NewLevelPanelEntry() *LevelPanelEntry {
	return &LevelPanelEntry{}
}

func (entry *LevelPanelEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
