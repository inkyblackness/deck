package ui

import (
	check "gopkg.in/check.v1"
)

type LimitedAnchorSuite struct {
}

var _ = check.Suite(&LimitedAnchorSuite{})

func (suite *LimitedAnchorSuite) TestValueReturnsReferenceValueIfWithinLimits(c *check.C) {
	from := NewAbsoluteAnchor(10.0)
	to := NewAbsoluteAnchor(20.0)
	reference := NewAbsoluteAnchor(15.0)
	anchor := NewLimitedAnchor(from, to, reference)

	c.Check(anchor.Value(), check.Equals, float32(15.0))
}

func (suite *LimitedAnchorSuite) TestValueReturnsFromValueIfReferenceSmallerThanFrom(c *check.C) {
	from := NewAbsoluteAnchor(10.0)
	to := NewAbsoluteAnchor(20.0)
	reference := NewAbsoluteAnchor(9.0)
	anchor := NewLimitedAnchor(from, to, reference)

	c.Check(anchor.Value(), check.Equals, float32(10.0))
}

func (suite *LimitedAnchorSuite) TestValueReturnsToValueIfReferenceGreaterThanTo(c *check.C) {
	from := NewAbsoluteAnchor(10.0)
	to := NewAbsoluteAnchor(20.0)
	reference := NewAbsoluteAnchor(21.0)
	anchor := NewLimitedAnchor(from, to, reference)

	c.Check(anchor.Value(), check.Equals, float32(20.0))
}

func (suite *LimitedAnchorSuite) TestValueReturnsCenterBetweenToAndFromIfReversed(c *check.C) {
	from := NewAbsoluteAnchor(30.0)
	to := NewAbsoluteAnchor(20.0)
	reference := NewAbsoluteAnchor(0.0)
	anchor := NewLimitedAnchor(from, to, reference)

	c.Check(anchor.Value(), check.Equals, float32(25.0))
}

func (suite *LimitedAnchorSuite) TestRequestValueUpdatesReferenceIfWithinLimits(c *check.C) {
	from := NewAbsoluteAnchor(30.0)
	to := NewAbsoluteAnchor(40.0)
	reference := NewAbsoluteAnchor(30.0)
	anchor := NewLimitedAnchor(from, to, reference)

	anchor.RequestValue(35.0)

	c.Check(reference.Value(), check.Equals, float32(35.0))
}

func (suite *LimitedAnchorSuite) TestRequestValueClipsReferenceAtToIfGreaterThanTo(c *check.C) {
	from := NewAbsoluteAnchor(30.0)
	to := NewAbsoluteAnchor(40.0)
	reference := NewAbsoluteAnchor(30.0)
	anchor := NewLimitedAnchor(from, to, reference)

	anchor.RequestValue(45.0)

	c.Check(anchor.Value(), check.Equals, float32(40.0))
}

func (suite *LimitedAnchorSuite) TestRequestValueClipsReferenceAtFromIfSmallerThanFrom(c *check.C) {
	from := NewAbsoluteAnchor(30.0)
	to := NewAbsoluteAnchor(40.0)
	reference := NewAbsoluteAnchor(30.0)
	anchor := NewLimitedAnchor(from, to, reference)

	anchor.RequestValue(25.0)

	c.Check(reference.Value(), check.Equals, float32(30.0))
}

func (suite *LimitedAnchorSuite) TestRequestValueIsIgnoredIfFromAndToReversed(c *check.C) {
	from := NewAbsoluteAnchor(50.0)
	to := NewAbsoluteAnchor(20.0)
	reference := NewAbsoluteAnchor(0.0)
	anchor := NewLimitedAnchor(from, to, reference)

	anchor.RequestValue(25.0)

	c.Check(reference.Value(), check.Equals, float32(0.0))
}
