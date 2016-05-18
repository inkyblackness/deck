package data

import (
	"fmt"
)

// LevelWeaponEntrySize specifies the byte count of a serialized LevelWeaponEntry.
const LevelWeaponEntrySize int = LevelObjectPrefixSize + 2

// LevelWeaponEntryCount specifies the count how many entries are in one level.
const LevelWeaponEntryCount int = 16

// LevelWeaponEntry describes a 'weapon' level object.
type LevelWeaponEntry struct {
	LevelObjectPrefix

	AmmoTypeOrCharge       byte
	AmmoCountOrTemperature byte
}

// NewLevelWeaponEntry returns a new instance of a LevelWeaponEntry.
func NewLevelWeaponEntry() *LevelWeaponEntry {
	return &LevelWeaponEntry{}
}

func (entry *LevelWeaponEntry) String() (result string) {
	result += entry.LevelObjectPrefix.String()
	result += fmt.Sprintf("Ammo Type/Charge: %d\n", entry.AmmoTypeOrCharge)
	result += fmt.Sprintf("Ammo Count/Temp: %d\n", entry.AmmoCountOrTemperature)

	return
}
