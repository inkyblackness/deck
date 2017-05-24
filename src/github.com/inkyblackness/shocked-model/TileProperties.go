package model

// TileProperties describe one tile in the map.
type TileProperties struct {
	Type *TileType

	FloorHeight   *HeightUnit
	CeilingHeight *HeightUnit
	SlopeHeight   *HeightUnit

	SlopeControl *SlopeControl

	CalculatedWallHeights *CalculatedWallHeights

	MusicIndex *int

	RealWorld *RealWorldTileProperties
	//Cyberspace *CyberspaceTileProperties
}
