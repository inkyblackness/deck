package data

// LevelProjectileEntrySize specifies the byte count of a serialized LevelProjectileEntry.
const LevelProjectileEntrySize int = LevelObjectPrefixSize + 34

// LevelProjectileEntryCount specifies the count how many entries are in one level.
const LevelProjectileEntryCount int = 32

// LevelProjectileEntry describes an 'Projectile' level object.
type LevelProjectileEntry struct {
	LevelObjectPrefix

	Unknown [34]byte
}

// NewLevelProjectileEntry returns a new LevelProjectileEntry instance.
func NewLevelProjectileEntry() *LevelProjectileEntry {
	return &LevelProjectileEntry{}
}

func (entry *LevelProjectileEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()

	return
}
