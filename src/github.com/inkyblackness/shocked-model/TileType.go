package model

// TileType describes the basic type of a tile.
type TileType string

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
