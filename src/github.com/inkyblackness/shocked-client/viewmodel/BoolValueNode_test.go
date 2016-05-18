package viewmodel

import (
	check "gopkg.in/check.v1"
)

type BoolValueNodeSuite struct {
}

var _ = check.Suite(&BoolValueNodeSuite{})

func (suite *BoolValueNodeSuite) TestSpecializeCallsBoolValue(c *check.C) {
	node := NewBoolValueNode("", false)
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.boolValueNodes, check.DeepEquals, []Node{node})
}

func (suite *BoolValueNodeSuite) TestLabel(c *check.C) {
	c.Check(NewBoolValueNode("l1", false).Label(), check.Equals, "l1")
	c.Check(NewBoolValueNode("l2", false).Label(), check.Equals, "l2")
}

func (suite *BoolValueNodeSuite) TestGetReturnsInitialValue(c *check.C) {
	c.Check(NewBoolValueNode("", true).Get(), check.Equals, true)
	c.Check(NewBoolValueNode("", false).Get(), check.Equals, false)
}

func (suite *BoolValueNodeSuite) TestSetChangesCurrentValue(c *check.C) {
	node := NewBoolValueNode("", true)

	node.Set(false)

	c.Check(node.Get(), check.Equals, false)
}

func (suite *BoolValueNodeSuite) TestSetCallsRegisteredSubscriberWithNewValue(c *check.C) {
	node := NewBoolValueNode("", false)
	capturedValue := false

	node.Subscribe(func(newValue bool) {
		capturedValue = newValue
	})
	node.Set(true)

	c.Check(capturedValue, check.Equals, true)
}
