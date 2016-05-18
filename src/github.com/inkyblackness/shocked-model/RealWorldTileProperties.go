package model

type RealWorldTileProperties struct {
	FloorTexture   *int `json:"floorTexture"`
	CeilingTexture *int `json:"ceilingTexture"`
	WallTexture    *int `json:"wallTexture"`

	FloorTextureRotations   *int `json:"floorTextureRotations"`
	CeilingTextureRotations *int `json:"ceilingTextureRotations"`

	UseAdjacentWallTexture *bool       `json:"useAdjacentWallTexture"`
	WallTextureOffset      *HeightUnit `json:"wallTextureOffset"`
}
