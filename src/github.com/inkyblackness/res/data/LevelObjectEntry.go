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

	X    MapCoordinate
	Y    MapCoordinate
	Z    byte
	Rot1 byte
	Rot2 byte
	Rot3 byte

	Unknown0013 [1]byte

	Type      res.ObjectType
	Hitpoints uint16

	Extra [4]byte
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
	result += fmt.Sprintf("Hitpoints: %d\n", entry.Hitpoints)

	return
}

// IsInUse returns true if the entry is active.
func (entry *LevelObjectEntry) IsInUse() bool {
	return entry.InUse != 0
}

// NextIndex returns the index of the next object.
func (entry *LevelObjectEntry) NextIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(entry.Next)
}

// SetNextIndex sets the index of the next object.
func (entry *LevelObjectEntry) SetNextIndex(index LevelObjectChainIndex) {
	entry.Next = uint16(index)
}

// PreviousIndex returns the index of the previous object.
func (entry *LevelObjectEntry) PreviousIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(entry.Previous)
}

// SetPreviousIndex sets the index of the previous object.
func (entry *LevelObjectEntry) SetPreviousIndex(index LevelObjectChainIndex) {
	entry.Previous = uint16(index)
}

// ReferenceIndex returns the index of a referenced object.
func (entry *LevelObjectEntry) ReferenceIndex() LevelObjectChainIndex {
	return LevelObjectChainIndex(entry.CrossReferenceTableIndex)
}

// SetReferenceIndex sets the index of a referenced object.
func (entry *LevelObjectEntry) SetReferenceIndex(index LevelObjectChainIndex) {
	entry.CrossReferenceTableIndex = uint16(index)
}
