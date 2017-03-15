package ui

import (
	check "gopkg.in/check.v1"
)

type ZeroAnchorSuite struct {
}

var _ = check.Suite(&ZeroAnchorSuite{})

func (suite *ZeroAnchorSuite) TestValueReturnsZero(c *check.C) {
	anchor := ZeroAnchor()

	c.Check(anchor.Value(), check.Equals, float32(0.0))
}

func (suite *ZeroAnchorSuite) TestRequestValueIsIgnored(c *check.C) {
	anchor := ZeroAnchor()

	anchor.RequestValue(123.0)

	c.Check(anchor.Value(), check.Equals, float32(0.0))
}
