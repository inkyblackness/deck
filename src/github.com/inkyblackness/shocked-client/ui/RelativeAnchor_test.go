package ui

import (
	check "gopkg.in/check.v1"
)

type RelativeAnchorSuite struct {
}

var _ = check.Suite(&RelativeAnchorSuite{})

func (suite *RelativeAnchorSuite) TestValueReturnsInitialValue_A(c *check.C) {
	from := NewAbsoluteAnchor(10.0)
	to := NewAbsoluteAnchor(20.0)
	anchor := NewRelativeAnchor(from, to, 0.5)

	c.Check(anchor.Value(), check.Equals, float32(15.0))
}

func (suite *RelativeAnchorSuite) TestValueReturnsInitialValue_B(c *check.C) {
	from := NewAbsoluteAnchor(30.0)
	to := NewAbsoluteAnchor(40.0)
	anchor := NewRelativeAnchor(from, to, 0.1)

	c.Check(anchor.Value(), check.Equals, float32(31.0))
}

func (suite *RelativeAnchorSuite) TestRequestValueUpdatesFraction_A(c *check.C) {
	from := NewAbsoluteAnchor(10.0)
	to := NewAbsoluteAnchor(20.0)
	anchor := NewRelativeAnchor(from, to, 0.6)

	anchor.RequestValue(12.0)
	from.RequestValue(5.0)
	to.RequestValue(10.0)

	c.Check(anchor.Value(), check.Equals, float32(6.0))
}

func (suite *RelativeAnchorSuite) TestRequestValueUpdatesFraction_B(c *check.C) {
	from := NewAbsoluteAnchor(0.0)
	to := NewAbsoluteAnchor(10.0)
	anchor := NewRelativeAnchor(from, to, 0.1)

	anchor.RequestValue(9.0)
	from.RequestValue(50.0)
	to.RequestValue(100.0)

	c.Check(anchor.Value(), check.Equals, float32(95.0))
}

func (suite *RelativeAnchorSuite) TestRequestValueIsIgnoredIfZeroDistance_B(c *check.C) {
	from := NewAbsoluteAnchor(0.0)
	to := NewAbsoluteAnchor(10.0)
	anchor := NewRelativeAnchor(from, to, 0.1)

	from.RequestValue(10.0)
	anchor.RequestValue(10.0)
	from.RequestValue(5.0)

	c.Check(anchor.Value(), check.Equals, float32(5.5))
}
