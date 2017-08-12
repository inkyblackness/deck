package data

import (
	check "gopkg.in/check.v1"
)

type BlockPuzzleStateSuite struct {
	data [16]byte
}

var _ = check.Suite(&BlockPuzzleStateSuite{})

func (suite *BlockPuzzleStateSuite) SetUpTest(c *check.C) {
	for index := 0; index < len(suite.data); index++ {
		suite.data[index] = 0x00
	}
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsTheCellValue_A(c *check.C) {
	state := NewBlockPuzzleState(suite.data[:], 1, 1)

	c.Check(state.CellValue(0, 0), check.Equals, 0)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsTheCellValue_B(c *check.C) {
	suite.data[16-4] = 2
	state := NewBlockPuzzleState(suite.data[:], 1, 1)

	c.Check(state.CellValue(0, 0), check.Equals, 2)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsTheCellValue_C(c *check.C) {
	suite.data[16-4] = 0x18
	state := NewBlockPuzzleState(suite.data[:], 2, 1)

	c.Check(state.CellValue(1, 0), check.Equals, 3)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsTheCellValue_D(c *check.C) {
	suite.data[15] = 0x40
	suite.data[8] = 0x01
	state := NewBlockPuzzleState(suite.data[:], 4, 3)

	c.Check(state.CellValue(3, 1), check.Equals, 5)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsTheCellValue_E(c *check.C) {
	suite.data[3] = 0x38
	state := NewBlockPuzzleState(suite.data[:], 6, 7)

	c.Check(state.CellValue(5, 6), check.Equals, 7)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsZeroOutOfBounds_A(c *check.C) {
	for index := 0; index < len(suite.data); index++ {
		suite.data[index] = 0xFF
	}
	state := NewBlockPuzzleState(suite.data[:], 1, 1)

	c.Check(state.CellValue(0, 1), check.Equals, 0)
	c.Check(state.CellValue(1, 0), check.Equals, 0)
	c.Check(state.CellValue(5, 5), check.Equals, 0)
	c.Check(state.CellValue(8, 8), check.Equals, 0)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsZeroOutOfBounds_B(c *check.C) {
	for index := 0; index < len(suite.data); index++ {
		suite.data[index] = 0xFF
	}
	state := NewBlockPuzzleState(suite.data[:], 7, 6)

	c.Check(state.CellValue(6, 7), check.Equals, 0)
	c.Check(state.CellValue(10, 10), check.Equals, 0)
}

func (suite *BlockPuzzleStateSuite) TestCellValueReturnsZeroOutOfBounds_C(c *check.C) {
	for index := 0; index < len(suite.data); index++ {
		suite.data[index] = 0xFF
	}
	state := NewBlockPuzzleState(suite.data[:], 7, 7)

	c.Check(state.CellValue(6, 0), check.Equals, 0)
}

func (suite *BlockPuzzleStateSuite) TestSetCellValue_A(c *check.C) {
	state := NewBlockPuzzleState(suite.data[:], 1, 1)

	state.SetCellValue(0, 0, 1)
	c.Check(suite.data, check.DeepEquals, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00})
}

func (suite *BlockPuzzleStateSuite) TestSetCellValue_B(c *check.C) {
	state := NewBlockPuzzleState(suite.data[:], 2, 2)

	state.SetCellValue(1, 0, 7)
	c.Check(suite.data, check.DeepEquals, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xC0, 0x01, 0x00, 0x00})
}

func (suite *BlockPuzzleStateSuite) TestSetCellValue_C(c *check.C) {
	state := NewBlockPuzzleState(suite.data[:], 6, 7)

	state.SetCellValue(5, 6, 5)
	c.Check(suite.data, check.DeepEquals, [16]byte{
		0x00, 0x00, 0x00, 0x28,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00})
}

func (suite *BlockPuzzleStateSuite) TestSetCellValueIgnoredOutOfBounds(c *check.C) {
	state := NewBlockPuzzleState(suite.data[:], 6, 7)

	state.SetCellValue(8, 8, 7)
	c.Check(suite.data, check.DeepEquals, [16]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00})
}
