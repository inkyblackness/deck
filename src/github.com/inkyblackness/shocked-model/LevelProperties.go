package model

// LevelProperties contains basic level information.
type LevelProperties struct {
	HeightShift    *int
	CyberspaceFlag *bool

	CeilingHasRadiation *bool
	CeilingEffectLevel  *int

	FloorHasBiohazard *bool
	FloorHasGravity   *bool
	FloorEffectLevel  *int
}
