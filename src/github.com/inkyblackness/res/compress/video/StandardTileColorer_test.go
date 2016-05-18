package video

import (
	check "gopkg.in/check.v1"
)

type StandardTileColorerSuite struct {
	buffer []byte
	stride int

	colorer TileColorFunction
}

var _ = check.Suite(&StandardTileColorerSuite{})

func (suite *StandardTileColorerSuite) SetUpTest(c *check.C) {
	suite.buffer = make([]byte, PixelPerTile*9)
	suite.stride = TileSideLength * 3

	for i := 0; i < len(suite.buffer); i++ {
		suite.buffer[i] = 0xDD
	}

	suite.colorer = StandardTileColorer(suite.buffer, suite.stride)
}

func (suite *StandardTileColorerSuite) TestFunctionColorsAllPixelOfATile(c *check.C) {
	suite.whenTileIsColored(1, 1, []byte{0xFF, 0x00}, 0, 1)

	suite.thenTileShouldBe(c, 1, 1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
}

func (suite *StandardTileColorerSuite) TestColoringATileLeavesPixelOfOthersAlone(c *check.C) {
	suite.whenTileIsColored(1, 1, []byte{0xFF}, 0, 1)

	suite.thenTileShouldBeUntouched(c, 0, 0)
	suite.thenTileShouldBeUntouched(c, 1, 0)
	suite.thenTileShouldBeUntouched(c, 2, 0)
	suite.thenTileShouldBeUntouched(c, 0, 1)
	//suite.thenTileShouldBeUntouched(c, 1, 1)
	suite.thenTileShouldBeUntouched(c, 2, 1)
	suite.thenTileShouldBeUntouched(c, 0, 2)
	suite.thenTileShouldBeUntouched(c, 1, 2)
	suite.thenTileShouldBeUntouched(c, 2, 2)
}

func (suite *StandardTileColorerSuite) TestFunctionWorksWithMaximumIndexSize(c *check.C) {
	suite.whenTileIsColored(1, 1, []byte{0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		0xFEDCBA9876543210, 4)

	suite.thenTileShouldBe(c, 1, 1, []byte{0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
}

func (suite *StandardTileColorerSuite) TestFunctionSkipsZeroPaletteIndices(c *check.C) {
	suite.givenTileIsFilled(1, 1, 0xFF)

	suite.whenTileIsColored(1, 1, []byte{0x33, 0x00}, 0xAAAA, 1)

	suite.thenTileShouldBe(c, 1, 1, []byte{0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF, 0x33, 0xFF})
}

func (suite *StandardTileColorerSuite) givenTileIsFilled(hTile int, vTile int, initValue byte) {
	suite.colorer(hTile, vTile, []byte{initValue, 0xEE}, 0x0000, 1)
}

func (suite *StandardTileColorerSuite) whenTileIsColored(hTile int, vTile int, lookupArray []byte, mask uint64, indexBitSize uint64) {
	suite.colorer(hTile, vTile, lookupArray, mask, indexBitSize)
}

func (suite *StandardTileColorerSuite) thenTileShouldBe(c *check.C, hTile, vTile int, expected []byte) {
	tile := suite.getTile(hTile, vTile)

	c.Check(tile, check.DeepEquals, expected)
}

func (suite *StandardTileColorerSuite) thenTileShouldBeUntouched(c *check.C, hTile, vTile int) {
	expected := make([]byte, PixelPerTile)
	for i := 0; i < len(expected); i++ {
		expected[i] = 0xDD
	}

	suite.thenTileShouldBe(c, hTile, vTile, expected)
}

func (suite *StandardTileColorerSuite) getTile(hTile, vTile int) []byte {
	tileBuffer := make([]byte, PixelPerTile)
	start := vTile*TileSideLength*suite.stride + hTile*TileSideLength

	for i := 0; i < TileSideLength; i++ {
		outOffset := TileSideLength * i
		inOffset := start + suite.stride*i
		copy(tileBuffer[outOffset:outOffset+TileSideLength], suite.buffer[inOffset:inOffset+TileSideLength])
	}

	return tileBuffer
}
