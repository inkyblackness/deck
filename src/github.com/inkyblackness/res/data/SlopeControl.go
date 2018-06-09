package data

// SlopeControl specifies how the floor and ceiling are affected for a sloped tile.
type SlopeControl uint32

// Constants for the slope controls.
const (
	SlopeCeilingInverted = SlopeControl(0)
	SlopeCeilingMirrored = SlopeControl(1)
	SlopeCeilingFlat     = SlopeControl(2)
	SlopeFloorFlat       = SlopeControl(3)
)
