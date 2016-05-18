package util

import (
	check "gopkg.in/check.v1"
)

type ValueUnifierSuite struct {
	store   *ValueUnifier
	queries map[int]int
}

var _ = check.Suite(&ValueUnifierSuite{})

func (suite *ValueUnifierSuite) TestValueReturnsDefaultForNoAddedValues(c *check.C) {
	unifier := NewValueUnifier("test")

	c.Check(unifier.Value(), check.Equals, "test")
}

func (suite *ValueUnifierSuite) TestWorksForIntegersAsWell(c *check.C) {
	unifier := NewValueUnifier(1234)

	c.Check(unifier.Value(), check.Equals, 1234)
}

func (suite *ValueUnifierSuite) TestValueReturnsUniqueValueIfAllAddedValuesWereTheSame_Integer(c *check.C) {
	unifier := NewValueUnifier(1234)

	unifier.Add(4)
	unifier.Add(4)
	unifier.Add(4)

	c.Check(unifier.Value(), check.Equals, 4)
}

func (suite *ValueUnifierSuite) TestValueReturnsUniqueValueIfAllAddedValuesWereTheSame_String(c *check.C) {
	unifier := NewValueUnifier("")

	unifier.Add("a")
	unifier.Add("a")
	unifier.Add("a")

	c.Check(unifier.Value(), check.Equals, "a")
}

func (suite *ValueUnifierSuite) TestValueReturnsDefaultForNonUniqueValues(c *check.C) {
	unifier := NewValueUnifier("def")

	unifier.Add("a")
	unifier.Add("b")
	unifier.Add("a")

	c.Check(unifier.Value(), check.Equals, "def")
}
