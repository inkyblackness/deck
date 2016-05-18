package data

// LevelAmmoEntrySize specifies the byte count of a serialized LevelAmmoEntry.
const LevelAmmoEntrySize int = LevelObjectPrefixSize + 0

// LevelAmmoEntryCount specifies the count how many entries are in one level.
const LevelAmmoEntryCount int = 32

// LevelAmmoEntry describes a 'Ammo' level object.
type LevelAmmoEntry struct {
	LevelObjectPrefix
}

// NewLevelAmmoEntry returns a new instance of a LevelAmmoEntry.
func NewLevelAmmoEntry() *LevelAmmoEntry {
	return &LevelAmmoEntry{}
}

func (entry *LevelAmmoEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
