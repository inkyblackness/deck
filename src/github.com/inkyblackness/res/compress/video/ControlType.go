package video

// ControlType defines how to interpret a ControlWord
type ControlType byte

const (
	CtrlColorTile2ColorsStatic  = ControlType(0)
	CtrlColorTile2ColorsMasked  = ControlType(1)
	CtrlColorTile4ColorsMasked  = ControlType(2)
	CtrlColorTile8ColorsMasked  = ControlType(3)
	CtrlColorTile16ColorsMasked = ControlType(4)

	CtrlSkip = ControlType(5)

	CtrlRepeatPrevious = ControlType(6)
	CtrlUnknown        = ControlType(7)
)
