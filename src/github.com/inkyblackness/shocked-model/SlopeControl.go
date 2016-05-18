package model

// SlopeControl specifies how the floor and ceiling are affected for a sloped tile.
type SlopeControl string

const (
	SlopeCeilingInverted = "ceilingInverted"
	SlopeCeilingMirrored = "ceilingMirrored"
	SlopeCeilingFlat     = "ceilingFlat"
	SlopeFloorFlat       = "floorFlat"
)
