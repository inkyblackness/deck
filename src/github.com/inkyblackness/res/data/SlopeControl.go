package data

// SlopeControl specifies how the floor and ceiling are affected for a sloped tile.
type SlopeControl uint32

const (
	SlopeCeilingInverted = SlopeControl(0)
	SlopeCeilingMirrored = SlopeControl(1)
	SlopeCeilingFlat     = SlopeControl(2)
	SlopeFloorFlat       = SlopeControl(3)
)
