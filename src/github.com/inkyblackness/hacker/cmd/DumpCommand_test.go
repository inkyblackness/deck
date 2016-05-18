package cmd

import (
	check "gopkg.in/check.v1"
)

type DumpCommandSuite struct {
	target *testTarget
}

var _ = check.Suite(&DumpCommandSuite{})

func (suite *DumpCommandSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
}

func (suite *DumpCommandSuite) TestDumpCommandReturnsNilForUnknownText(c *check.C) {
	result := dumpCommand("not dump")

	c.Assert(result, check.IsNil)
}

func (suite *DumpCommandSuite) TestDumpCommandReturnsDumpCall(c *check.C) {
	result := dumpCommand("dump")

	result(suite.target)

	c.Assert(len(suite.target.dumpParam), check.Equals, 1)
}
