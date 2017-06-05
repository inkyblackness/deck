package model

// CyberspaceTileProperties describes tile properties of cyberspace.
type CyberspaceTileProperties struct {
	FloorColorIndex   *int
	CeilingColorIndex *int

	FlightPullType *int

	GameOfLifeSet *bool
}
