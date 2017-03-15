package ui

import (
	check "gopkg.in/check.v1"
)

type ResolvingAnchorSuite struct {
}

var _ = check.Suite(&ResolvingAnchorSuite{})

func (suite *ResolvingAnchorSuite) TestValueReturnsInitialValue_A(c *check.C) {
	referred := NewAbsoluteAnchor(20.0)
	anchor := NewResolvingAnchor(func() Anchor { return referred })

	c.Check(anchor.Value(), check.Equals, float32(20.0))
}

func (suite *ResolvingAnchorSuite) TestValueReturnsInitialValue_B(c *check.C) {
	referred := NewAbsoluteAnchor(10.0)
	anchor := NewResolvingAnchor(func() Anchor { return referred })

	c.Check(anchor.Value(), check.Equals, float32(10.0))
}

func (suite *ResolvingAnchorSuite) TestRequestValueForwardsValue_A(c *check.C) {
	referred := NewAbsoluteAnchor(10.0)
	anchor := NewResolvingAnchor(func() Anchor { return referred })

	anchor.RequestValue(30.0)

	c.Check(referred.Value(), check.Equals, float32(30.0))
}

func (suite *ResolvingAnchorSuite) TestRequestValueUpdatesValue_B(c *check.C) {
	referred := NewAbsoluteAnchor(20.0)
	anchor := NewResolvingAnchor(func() Anchor { return referred })

	anchor.RequestValue(35.0)

	c.Check(referred.Value(), check.Equals, float32(35.0))
}

func (suite *ResolvingAnchorSuite) TestChangesOfReferredAreUpdated(c *check.C) {
	referred := ZeroAnchor()
	anchor := NewResolvingAnchor(func() Anchor { return referred })

	referred = NewAbsoluteAnchor(10.0)

	c.Check(anchor.Value(), check.Equals, float32(10.0))
}

func (suite *ResolvingAnchorSuite) TestReferredOfNilDefaultsToZeroAnchor_Value(c *check.C) {
	anchor := NewResolvingAnchor(func() Anchor { return nil })

	c.Check(anchor.Value(), check.Equals, float32(0.0))
}

func (suite *ResolvingAnchorSuite) TestReferredOfNilDefaultsToZeroAnchor_RequestValue(c *check.C) {
	anchor := NewResolvingAnchor(func() Anchor { return nil })

	anchor.RequestValue(10.0)

	c.Check(anchor.Value(), check.Equals, float32(0.0))
}
