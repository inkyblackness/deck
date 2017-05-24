package model

// RealWorldTileProperties describes tile properties of the real wold.
type RealWorldTileProperties struct {
	FloorTexture   *int
	CeilingTexture *int
	WallTexture    *int

	FloorTextureRotations   *int
	CeilingTextureRotations *int

	UseAdjacentWallTexture *bool
	WallTextureOffset      *HeightUnit

	FloorHazard   *bool
	CeilingHazard *bool

	FloorShadow   *int
	CeilingShadow *int
}
