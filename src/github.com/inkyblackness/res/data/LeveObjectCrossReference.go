package data

import (
	"fmt"
)

// LevelObjectCrossReferenceSize specifies the byte count of a serialized LevelObjectCrossReference.
const LevelObjectCrossReferenceSize int = 10

// LevelObjectCrossReference describes the cross reference between one or more map tiles and level objects.
type LevelObjectCrossReference struct {
	TileX uint16
	TileY uint16

	LevelObjectTableIndex uint16

	NextObjectIndex uint16
	NextTileIndex   uint16
}

// DefaultLevelObjectCrossReference returns a new instance of LevelObjectCrossReference.
func DefaultLevelObjectCrossReference() *LevelObjectCrossReference {
	return &LevelObjectCrossReference{}
}

func (ref *LevelObjectCrossReference) String() (result string) {
	result += fmt.Sprintf("Coord: X: %v Y: %v\n", ref.TileX, ref.TileY)
	result += fmt.Sprintf("Level Object Table Index: %d\n", ref.LevelObjectTableIndex)
	result += fmt.Sprintf("Next Object Index: %d\n", ref.NextObjectIndex)
	result += fmt.Sprintf("Next Tile Index: %d\n", ref.NextTileIndex)

	return
}
