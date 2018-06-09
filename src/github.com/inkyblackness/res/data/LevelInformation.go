package data

import (
	"fmt"
)

const (
	defaultMapDimension      = 64
	defaultMapDimensionShift = 6
	defaultHeightShift       = 3
	defaultTimerValue1       = 64
	defaultTimerValue2       = 8
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

	Unknown001C [9]byte

	TimerValue1 byte
	Unknown0026 [3]byte
	TimerCount  byte
	Unknown002A [3]byte
	TimerValue2 byte
	Unknown002E [12]byte
}

// DefaultLevelInformation returns an instance of LevelInformation with default values.
func DefaultLevelInformation() *LevelInformation {
	info := &LevelInformation{
		Unknown0000: defaultMapDimension,
		Unknown0004: defaultMapDimension,
		Unknown0008: defaultMapDimensionShift,
		Unknown000C: defaultMapDimensionShift,
		HeightShift: defaultHeightShift,
		TimerValue1: defaultTimerValue1,
		TimerValue2: defaultTimerValue2}

	return info
}

func (info *LevelInformation) String() (result string) {
	result += fmt.Sprintf("Cyberspace: %v\n", info.IsCyberspace())
	result += fmt.Sprintf("Height Shift: %d\n", info.HeightShift)
	result += fmt.Sprintf("Timer Count: %d\n", info.TimerCount)

	return result
}

// IsCyberspace returns true for cyberspace levels
func (info *LevelInformation) IsCyberspace() bool {
	return info.CyberspaceFlag != 0
}
