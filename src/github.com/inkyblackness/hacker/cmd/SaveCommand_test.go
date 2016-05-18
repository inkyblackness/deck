package cmd

import (
	check "gopkg.in/check.v1"
)

type SaveCommandSuite struct {
	target *testTarget
}

var _ = check.Suite(&SaveCommandSuite{})

func (suite *SaveCommandSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
}

func (suite *SaveCommandSuite) TestSaveCommandReturnsNilForUnknownText(c *check.C) {
	result := saveCommand("not save")

	c.Assert(result, check.IsNil)
}

func (suite *SaveCommandSuite) TestSaveCommandReturnsCallsInfo(c *check.C) {
	result := saveCommand("save")

	result(suite.target)

	c.Assert(len(suite.target.saveParam), check.Equals, 1)
}
