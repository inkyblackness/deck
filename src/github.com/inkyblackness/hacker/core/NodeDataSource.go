package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// NodeDataSource implements the query.DataSource interface to provide
// archive data based on DataNodes
type NodeDataSource struct {
	archiveNode  DataNode
	hacker       *Hacker
	currentLevel int
}

// NewNodeDataSource returns a new instance of NodeDataSource.
func NewNodeDataSource(archiveNode DataNode, hacker *Hacker) *NodeDataSource {
	source := &NodeDataSource{
		archiveNode: archiveNode,
		hacker:      hacker}

	source.currentLevel = int(source.GameState().CurrentLevel)

	return source
}

func (source *NodeDataSource) levelChunkPath(level int, chunkID int) string {
	return fmt.Sprintf("%04X/0", 4000+100*level+chunkID)
}

func (source *NodeDataSource) levelEntryPath(level int, table int, index int) string {
	return fmt.Sprintf("%04X/0/%d", 4000+100*level+table, index)
}

func (source *NodeDataSource) currentLevelEntryPath(table int, index int) string {
	return source.levelEntryPath(source.currentLevel, table, index)
}

// ObjectEntryPath returns the path to the object class entry
func (source *NodeDataSource) ObjectEntryPath(class int, index int) string {
	return source.currentLevelEntryPath(10+class, index)
}

func (source *NodeDataSource) mapData(data interface{}, path string) {
	blockNode := source.hacker.resolveFrom(source.archiveNode, path)
	coder := serial.NewDecoder(bytes.NewReader(blockNode.Data()))
	serial.MapData(data, coder)
}

// GameState returns a GameState instance mapped to the archive data.
func (source *NodeDataSource) GameState() *data.GameState {
	gameState := data.DefaultGameState()

	source.mapData(gameState, "0FA1/0")

	return gameState
}

// Tile returns a tile map entry mapped to the current level.
func (source *NodeDataSource) Tile(x int, y int) *data.TileMapEntry {
	entry := data.DefaultTileMapEntry()
	path := source.currentLevelEntryPath(5, y*64+x)

	source.mapData(entry, path)

	return entry
}

// LevelObjectCrossReference returns a cross reference from given index mapped to the current level.
func (source *NodeDataSource) LevelObjectCrossReference(index uint16) *data.LevelObjectCrossReference {
	ref := data.DefaultLevelObjectCrossReference()

	source.mapData(ref, source.currentLevelEntryPath(9, int(index)))

	return ref
}

// LevelObject returns a level object entry from the current level
func (source *NodeDataSource) LevelObject(index uint16) *data.LevelObjectEntry {
	entry := data.DefaultLevelObjectEntry()

	source.mapData(entry, source.currentLevelEntryPath(8, int(index)))

	return entry
}

// LevelIDs returns the list of available levels
func (source *NodeDataSource) LevelIDs() (levelIDs []int) {
	for i := 0; i < 16; i++ {
		path := source.levelChunkPath(i, 2)
		resolved := source.hacker.resolveFrom(source.archiveNode, path)

		if resolved != nil {
			levelIDs = append(levelIDs, i)
		}
	}

	return
}

// LevelChunkData returns the raw data from the given chunk.
func (source *NodeDataSource) LevelChunkData(level int, levelChunk int) (data []byte) {
	resolved := source.hacker.resolveFrom(source.archiveNode, source.levelChunkPath(level, levelChunk))
	if resolved != nil {
		data = resolved.Data()
	}

	return
}
