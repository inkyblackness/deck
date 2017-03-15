package ui

import (
	check "gopkg.in/check.v1"
)

type OffsetAnchorSuite struct {
}

var _ = check.Suite(&OffsetAnchorSuite{})

func (suite *OffsetAnchorSuite) TestValueReturnsInitialValue_A(c *check.C) {
	base := NewAbsoluteAnchor(10.0)
	anchor := NewOffsetAnchor(base, 5.0)

	c.Check(anchor.Value(), check.Equals, float32(15.0))
}

func (suite *OffsetAnchorSuite) TestValueReturnsInitialValue_B(c *check.C) {
	base := NewAbsoluteAnchor(30.0)
	anchor := NewOffsetAnchor(base, -5.0)

	c.Check(anchor.Value(), check.Equals, float32(25.0))
}

func (suite *OffsetAnchorSuite) TestRequestValueUpdatesValue(c *check.C) {
	base := NewAbsoluteAnchor(0.0)
	anchor := NewOffsetAnchor(base, 10.0)

	anchor.RequestValue(15.0)

	c.Check(anchor.Value(), check.Equals, float32(15.0))
}

func (suite *OffsetAnchorSuite) TestRequestValueModifiesOffset_A(c *check.C) {
	base := NewAbsoluteAnchor(0.0)
	anchor := NewOffsetAnchor(base, 10.0)

	anchor.RequestValue(15.0)
	base.RequestValue(100.0)

	c.Check(anchor.Value(), check.Equals, float32(115.0))
}

func (suite *OffsetAnchorSuite) TestRequestValueModifiesOffset_B(c *check.C) {
	base := NewAbsoluteAnchor(50.0)
	anchor := NewOffsetAnchor(base, 10.0)

	anchor.RequestValue(70.0)
	base.RequestValue(10.0)

	c.Check(anchor.Value(), check.Equals, float32(30.0))
}
