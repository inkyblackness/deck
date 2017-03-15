package data

// TimerEntrySize is the number of bytes for one timer entry
const TimerEntrySize int = 8

// TimerEntry describes one timer of a level.
type TimerEntry struct {
	TriggerTime  uint16
	Unknown0002  uint16
	TargetObject uint16
	Unknown0006  uint16
}
