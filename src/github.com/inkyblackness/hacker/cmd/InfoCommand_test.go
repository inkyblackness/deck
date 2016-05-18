package cmd

import (
	check "gopkg.in/check.v1"
)

type InfoCommandSuite struct {
	target *testTarget
}

var _ = check.Suite(&InfoCommandSuite{})

func (suite *InfoCommandSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
}

func (suite *InfoCommandSuite) TestInfoCommandReturnsNilForUnknownText(c *check.C) {
	result := infoCommand("not load")

	c.Assert(result, check.IsNil)
}

func (suite *InfoCommandSuite) TestInfoCommandReturnsCallsInfo(c *check.C) {
	result := infoCommand("info")

	result(suite.target)

	c.Assert(len(suite.target.infoParam), check.Equals, 1)
}
