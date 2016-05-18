package data

import (
	"fmt"
)

// LevelVariablesSize specifies the byte count of a serialized LevelVariables.
const LevelVariablesSize int = 0x5E

// LevelVariables contains various variables about the level.
type LevelVariables struct {
	Size uint32

	Radiation         byte
	BioOrGravity      byte
	GravitySwitch     byte
	RadiationRegister byte
	BioRegister       byte

	Unknown0009 [85]byte
}

// NewLevelVariables returns a new instance of level variables.
func NewLevelVariables() *LevelVariables {
	return &LevelVariables{Size: uint32(LevelVariablesSize)}
}

func (info *LevelVariables) String() (result string) {
	result += fmt.Sprintf("Radiation: %f LBP\n", float32(info.Radiation)*0.5)
	result += fmt.Sprintf("Bio/Gravity: %f LBP / %d%% Gravity\n", float32(info.BioOrGravity)*0.5, int(info.BioOrGravity)*25)
	result += fmt.Sprintf("GravitySwitch: %v", info.IsGravityModified())

	return
}

// IsGravityModified returns true if gravity is non-standard 100%.
func (info *LevelVariables) IsGravityModified() bool {
	return info.GravitySwitch != 0
}
