package textprop

import (
	"fmt"
)

// Entry describes one texture properties entry
type Entry struct {
	IndexLowByte1       byte
	IndexLowByte2       byte
	Unknown0002         uint32
	Climbable           byte
	Unknown0007         [1]byte
	TransparencyControl byte
	AnimationGroup      byte
	AnimationIndex      byte
}

func (entry *Entry) String() (result string) {
	result += fmt.Sprintf("Index Bytes: %d/%d\n", entry.IndexLowByte1, entry.IndexLowByte2)
	result += fmt.Sprintf("Climbable: %v\n", entry.IsClimbable())
	result += fmt.Sprintf("TransparencyControl: 0x%02X\n", entry.TransparencyControl)
	result += fmt.Sprintf("Animation: %d:%d\n", entry.AnimationGroup, entry.AnimationIndex)

	return
}

func (entry *Entry) IsClimbable() bool {
	return entry.Climbable != 0
}
