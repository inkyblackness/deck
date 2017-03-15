package ui

import (
	check "gopkg.in/check.v1"
)

type AbsoluteAnchorSuite struct {
}

var _ = check.Suite(&AbsoluteAnchorSuite{})

func (suite *AbsoluteAnchorSuite) TestValueReturnsInitialValue_A(c *check.C) {
	anchor := NewAbsoluteAnchor(10.0)

	c.Check(anchor.Value(), check.Equals, float32(10.0))
}

func (suite *AbsoluteAnchorSuite) TestValueReturnsInitialValue_B(c *check.C) {
	anchor := NewAbsoluteAnchor(20.0)

	c.Check(anchor.Value(), check.Equals, float32(20.0))
}

func (suite *AbsoluteAnchorSuite) TestRequestValueUpdatesValue_A(c *check.C) {
	anchor := NewAbsoluteAnchor(20.0)

	anchor.RequestValue(30.0)

	c.Check(anchor.Value(), check.Equals, float32(30.0))
}

func (suite *AbsoluteAnchorSuite) TestRequestValueUpdatesValue_B(c *check.C) {
	anchor := NewAbsoluteAnchor(20.0)

	anchor.RequestValue(35.0)

	c.Check(anchor.Value(), check.Equals, float32(35.0))
}
