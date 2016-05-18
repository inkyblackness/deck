package cmd

import (
	check "gopkg.in/check.v1"
)

type PutCommandSuite struct {
	target *testTarget
}

var _ = check.Suite(&PutCommandSuite{})

func (suite *PutCommandSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
}

func (suite *PutCommandSuite) TestPutCommandReturnsNilForUnknownText(c *check.C) {
	result := putCommand("not put")

	c.Assert(result, check.IsNil)
}

func (suite *PutCommandSuite) TestPutRecognizesOffset(c *check.C) {
	result := putCommand("put 0 00")

	result(suite.target)

	c.Assert(len(suite.target.putParam), check.Equals, 1)
}

func (suite *PutCommandSuite) TestPutRecognizesOffsetAsHex(c *check.C) {
	result := putCommand("put 12AB 00")

	result(suite.target)

	c.Assert(suite.target.putParam[0][0], check.Equals, uint32(0x12AB))
}

func (suite *PutCommandSuite) TestPutRecognizesOffsetAsHexWithLeadingZeroes(c *check.C) {
	result := putCommand("put 000A 00")

	result(suite.target)

	c.Assert(suite.target.putParam[0][0], check.Equals, uint32(0x000A))
}

func (suite *PutCommandSuite) TestPutRecognizesBytes(c *check.C) {
	result := putCommand("put 12AB 00 0A 0c")

	result(suite.target)

	c.Assert(suite.target.putParam[0][1], check.DeepEquals, []byte{0x00, 0x0A, 0x0C})
}
