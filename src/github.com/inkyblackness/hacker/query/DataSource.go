package query

import (
	"github.com/inkyblackness/res/data"
)

// DataSource is an abstracting interface about accessing archive data
type DataSource interface {
	GameState() *data.GameState
	Tile(x int, y int) *data.TileMapEntry
	LevelObjectCrossReference(index uint16) *data.LevelObjectCrossReference
	LevelObject(index uint16) *data.LevelObjectEntry

	ObjectEntryPath(class int, index int) string

	LevelIDs() []int
	LevelChunkData(level int, chunk int) []byte
}
