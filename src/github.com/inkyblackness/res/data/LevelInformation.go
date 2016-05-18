package data

import (
	"fmt"
)

const (
	defaultMapDimension      uint32 = 64
	defaultMapDimensionShift uint32 = 6
	defaultHeightShift       uint32 = 3
)

// LevelInformation contains information about a single level.
type LevelInformation struct {
	Unknown0000 uint32
	Unknown0004 uint32
	Unknown0008 uint32
	Unknown000C uint32

	HeightShift uint32

	IgnoredPlaceholder uint32

	CyberspaceFlag uint32

	Unknown001C [30]byte
}

// DefaultLevelInformation returns an instance of LevelInformation with default values.
func DefaultLevelInformation() *LevelInformation {
	info := &LevelInformation{
		Unknown0000: defaultMapDimension,
		Unknown0004: defaultMapDimension,
		Unknown0008: defaultMapDimensionShift,
		Unknown000C: defaultMapDimensionShift,
		HeightShift: defaultHeightShift}

	return info
}

func (info *LevelInformation) String() (result string) {
	result += fmt.Sprintf("Cyberspace: %v\n", info.IsCyberspace())
	result += fmt.Sprintf("Height Shift: %d\n", info.HeightShift)

	return result
}

// IsCyberspace returns true for cyberspace levels
func (info *LevelInformation) IsCyberspace() bool {
	return info.CyberspaceFlag != 0
}
