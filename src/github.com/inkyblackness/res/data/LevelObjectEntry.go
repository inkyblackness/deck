package data

import (
	"fmt"

	"github.com/inkyblackness/res"
)

// LevelObjectEntrySize specifies the byte count of a serialized LevelObjectEntry.
const LevelObjectEntrySize int = 27

// LevelObjectEntry describes the basic information about a level object.
type LevelObjectEntry struct {
	InUse    byte
	Class    res.ObjectClass
	Subclass res.ObjectSubclass

	ClassTableIndex          uint16
	CrossReferenceTableIndex uint16
	Previous                 uint16
	Next                     uint16

	X    TileCoordinate
	Y    TileCoordinate
	Z    byte
	Rot1 byte
	Rot2 byte
	Rot3 byte

	Unknown0013 [1]byte

	Type res.ObjectType

	Unknown0015 [2]byte

	Unknown0017 [4]byte
}

// DefaultLevelObjectEntry returns a new LevelObjectEntry instance.
func DefaultLevelObjectEntry() *LevelObjectEntry {
	return &LevelObjectEntry{}
}

func (entry *LevelObjectEntry) String() (result string) {
	result += fmt.Sprintf("In Use: %v\n", entry.IsInUse())
	result += fmt.Sprintf("ObjectID: %d/%d/%d\n", entry.Class, entry.Subclass, entry.Type)
	result += fmt.Sprintf("Coord: X: %v Y: %v Z: %d\n", entry.X, entry.Y, entry.Z)
	result += fmt.Sprintf("Rotation: %d, %d, %d\n", entry.Rot1, entry.Rot2, entry.Rot3)
	result += fmt.Sprintf("Class Table Index: %d\n", entry.ClassTableIndex)
	result += fmt.Sprintf("Cross Reference Index: %d\n", entry.CrossReferenceTableIndex)
	result += fmt.Sprintf("Links: <- %d | %d ->\n", entry.Previous, entry.Next)

	return
}

// IsInUse returns true if the entry is active.
func (entry *LevelObjectEntry) IsInUse() bool {
	return entry.InUse != 0
}
