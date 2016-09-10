package logic

import (
	"github.com/inkyblackness/res/data"

	check "gopkg.in/check.v1"
)

type TileMapSuite struct {
}

var _ = check.Suite(&TileMapSuite{})

func (suite *TileMapSuite) SetUpTest(c *check.C) {
}

func (suite *TileMapSuite) TestDecodeTileMapRestoresMap(c *check.C) {
	width := 2
	height := 3
	tileMap := NewTileMap(width, height)

	tileMap.Entry(AtTile(1, 2)).Type = data.DiagonalOpenNorthEast
	serialized := tileMap.Encode()

	newMap := DecodeTileMap(serialized, width, height)

	c.Check(newMap.Entry(AtTile(1, 2)).Type, check.Equals, data.DiagonalOpenNorthEast)
}

func (suite *TileMapSuite) TestTileMapImplementsTileMapReferencer(c *check.C) {
	width := 7
	tileMap := NewTileMap(width, 5)

	location := AtTile(1, 3)
	tileMap.SetReferenceIndex(location, CrossReferenceListIndex(123))
	index := tileMap.ReferenceIndex(location)

	c.Check(index, check.Equals, CrossReferenceListIndex(123))
	c.Check(tileMap.tiles[3*width+1].FirstObjectIndex, check.Equals, uint16(123))
}
