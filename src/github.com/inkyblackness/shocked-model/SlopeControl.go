package model

// SlopeControl specifies how the floor and ceiling are affected for a sloped tile.
type SlopeControl string

// All known slope control values, as string.
const (
	SlopeCeilingInverted = "ceilingInverted"
	SlopeCeilingMirrored = "ceilingMirrored"
	SlopeCeilingFlat     = "ceilingFlat"
	SlopeFloorFlat       = "floorFlat"
)

// SlopeControls returns all available slope control values.
func SlopeControls() []SlopeControl {
	return []SlopeControl{SlopeCeilingInverted, SlopeCeilingMirrored, SlopeCeilingFlat, SlopeFloorFlat}
}
