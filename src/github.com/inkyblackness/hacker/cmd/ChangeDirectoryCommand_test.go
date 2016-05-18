package cmd

import (
	check "gopkg.in/check.v1"
)

type ChangeDirectoryCommandSuite struct {
	target *testTarget
}

var _ = check.Suite(&ChangeDirectoryCommandSuite{})

func (suite *ChangeDirectoryCommandSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
}

func (suite *ChangeDirectoryCommandSuite) TestChangeDirectoryCommandReturnsNilForUnknownText(c *check.C) {
	result := changeDirectoryCommand("not cd")

	c.Assert(result, check.IsNil)
}

func (suite *ChangeDirectoryCommandSuite) TestChangeDirectoryCommandReturnsFunctionForProperParameter(c *check.C) {
	result := changeDirectoryCommand("cd test")

	c.Assert(result, check.Not(check.IsNil))
}
func (suite *ChangeDirectoryCommandSuite) TestChangeDirectoryCommandReturnsCdFunctionForProperParameter(c *check.C) {
	result := changeDirectoryCommand("cd test")

	result(suite.target)

	c.Assert(len(suite.target.cdParam), check.Equals, 1)
}

func (suite *ChangeDirectoryCommandSuite) TestChangeDirectoryHandlesVariousParameterFormats(c *check.C) {
	suite.verifyCdParameter(c, `cd test`, "test")
	suite.verifyCdParameter(c, `cd   test`, "test")
	suite.verifyCdParameter(c, `cd test/test2`, "test/test2")
	suite.verifyCdParameter(c, `cd /one/2/three/..`, "/one/2/three/..")
}

func (suite *ChangeDirectoryCommandSuite) verifyCdParameter(c *check.C, input, path string) {
	result := changeDirectoryCommand(input)

	c.Assert(result, check.Not(check.IsNil))
	result(suite.target)

	c.Check(suite.target.cdParam[len(suite.target.cdParam)-1], check.DeepEquals, []interface{}{path})
}
