package cmd

import (
	"github.com/inkyblackness/hacker/styling"

	check "gopkg.in/check.v1"
)

type EvaluaterSuite struct {
	target *testTarget
	eval   *Evaluater
}

var _ = check.Suite(&EvaluaterSuite{})

func (suite *EvaluaterSuite) SetUpTest(c *check.C) {
	suite.target = &testTarget{}
	suite.eval = NewEvaluater(styling.NullStyle(), suite.target)
}

func (suite *EvaluaterSuite) TestEvaluateReturnsUnknownCommand(c *check.C) {
	result := suite.eval.Evaluate("dummy text")

	c.Assert(result, check.Equals, "Unknown command: [dummy text]")
}

func (suite *EvaluaterSuite) TestEvaluateUnderstandsCommands(c *check.C) {
	suite.verifyCommand(c, `load "a" "b"`, `Load("a", "b")`)
	suite.verifyCommand(c, `info`, `Info()`)
	suite.verifyCommand(c, `cd test`, `Cd(test)`)
	suite.verifyCommand(c, `dump`, `Dump()`)
	suite.verifyCommand(c, `save`, `Save()`)
	suite.verifyCommand(c, `put 0 01`, `Put(0, [1])`)
	suite.verifyCommand(c, `query local`, `Query(local)`)
}

func (suite *EvaluaterSuite) verifyCommand(c *check.C, input string, output string) {
	result := suite.eval.Evaluate(input)

	c.Check(result, check.Equals, output)
}
