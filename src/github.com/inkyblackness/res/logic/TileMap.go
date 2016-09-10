package logic

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res/data"
)

// TileMap contains the tiles of a map
type TileMap struct {
	width, height int
	tiles         []data.TileMapEntry
}

// NewTileMap returns a new instance of a tile map.
func NewTileMap(width, height int) *TileMap {
	tileMap := &TileMap{
		width:  width,
		height: height,
		tiles:  make([]data.TileMapEntry, width*height)}

	return tileMap
}

// DecodeTileMap decodes a tile map from a byte array. It expects the given
// size.
func DecodeTileMap(serialized []byte, width, height int) *TileMap {
	tileMap := NewTileMap(width, height)
	buf := bytes.NewReader(serialized)

	binary.Read(buf, binary.LittleEndian, &tileMap.tiles)

	return tileMap
}

// Encode serializes the map into a byte array.
func (tileMap *TileMap) Encode() []byte {
	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, &tileMap.tiles)

	return buf.Bytes()
}

// Entry returns the tile at the given location.
func (tileMap *TileMap) Entry(location TileLocation) *data.TileMapEntry {
	x, y := location.XY()

	return &tileMap.tiles[int(y)*tileMap.width+int(x)]
}

// ReferenceIndex returns the index of the first cross-reference for the given tile.
func (tileMap *TileMap) ReferenceIndex(location TileLocation) CrossReferenceListIndex {
	return CrossReferenceListIndex(tileMap.Entry(location).FirstObjectIndex)
}

// SetReferenceIndex sets the index of the first cross-reference for the given tile.
func (tileMap *TileMap) SetReferenceIndex(location TileLocation, index CrossReferenceListIndex) {
	tileMap.Entry(location).FirstObjectIndex = uint16(index)
}
