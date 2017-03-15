package model

// TileType describes the basic type of a tile.
type TileType string

// All known tile types, as string
const (
	Solid TileType = "solid"
	Open           = "open"

	DiagonalOpenSouthEast = "diagonalOpenSouthEast"
	DiagonalOpenSouthWest = "diagonalOpenSouthWest"
	DiagonalOpenNorthWest = "diagonalOpenNorthWest"
	DiagonalOpenNorthEast = "diagonalOpenNorthEast"

	SlopeSouthToNorth = "slopeSouthToNorth"
	SlopeWestToEast   = "slopeWestToEast"
	SlopeNorthToSouth = "slopeNorthToSouth"
	SlopeEastToWest   = "slopeEastToWest"

	ValleySouthEastToNorthWest = "valleySouthEastToNorthWest"
	ValleySouthWestToNorthEast = "valleySouthWestToNorthEast"
	ValleyNorthWestToSouthEast = "valleyNorthWestToSouthEast"
	ValleyNorthEastToSouthWest = "valleyNorthEastToSouthWest"

	RidgeNorthWestToSouthEast = "ridgeNorthWestToSouthEast"
	RidgeNorthEastToSouthWest = "ridgeNorthEastToSouthWest"
	RidgeSouthEastToNorthWest = "ridgeSouthEastToNorthWest"
	RidgeSouthWestToNorthEast = "ridgeSouthWestToNorthEast"
)

// TileTypes returns a list of all types.
func TileTypes() []TileType {
	return []TileType{
		Solid, Open,
		DiagonalOpenSouthEast, DiagonalOpenSouthWest, DiagonalOpenNorthWest, DiagonalOpenNorthEast,
		SlopeSouthToNorth, SlopeWestToEast, SlopeNorthToSouth, SlopeEastToWest,
		ValleySouthEastToNorthWest, ValleySouthWestToNorthEast, ValleyNorthWestToSouthEast, ValleyNorthEastToSouthWest,
		RidgeNorthWestToSouthEast, RidgeNorthEastToSouthWest, RidgeSouthEastToNorthWest, RidgeSouthWestToNorthEast}
}
