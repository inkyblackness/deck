package model

import (
	"fmt"

	check "gopkg.in/check.v1"
)

type TileCoordinateSuite struct {
	store   *TileCoordinate
	queries map[int]int
}

var _ = check.Suite(&TileCoordinateSuite{})

func (suite *TileCoordinateSuite) TestXYReturnsValuesA(c *check.C) {
	coord := TileCoordinateOf(10, 20)
	x, y := coord.XY()

	c.Check(x, check.Equals, 10)
	c.Check(y, check.Equals, 20)
}

func (suite *TileCoordinateSuite) TestXYReturnsValuesB(c *check.C) {
	coord := TileCoordinateOf(30, 40)
	x, y := coord.XY()

	c.Check(x, check.Equals, 30)
	c.Check(y, check.Equals, 40)
}

func (suite *TileCoordinateSuite) TestStringImplementsStringerInterface(c *check.C) {
	coord := TileCoordinateOf(40, 50)

	text := fmt.Sprintf("%v", coord)

	c.Check(text, check.Equals, "40/50")
}
