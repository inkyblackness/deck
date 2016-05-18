package viewmodel

import (
	check "gopkg.in/check.v1"
)

type StringValueNodeSuite struct {
}

var _ = check.Suite(&StringValueNodeSuite{})

func (suite *StringValueNodeSuite) TestSpecializeCallsStringValue(c *check.C) {
	node := NewStringValueNode("label", "abc")
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.stringValueNodes, check.DeepEquals, []Node{node})
}

func (suite *StringValueNodeSuite) TestLabel(c *check.C) {
	c.Check(NewStringValueNode("l1", "").Label(), check.Equals, "l1")
	c.Check(NewStringValueNode("l2", "").Label(), check.Equals, "l2")
}

func (suite *StringValueNodeSuite) TestGetReturnsInitialValue(c *check.C) {
	c.Check(NewStringValueNode("", "test").Get(), check.Equals, "test")
	c.Check(NewStringValueNode("", "other").Get(), check.Equals, "other")
}

func (suite *StringValueNodeSuite) TestSetChangesCurrentValue(c *check.C) {
	node := NewStringValueNode("", "first")

	node.Set("second")

	c.Check(node.Get(), check.Equals, "second")
}

func (suite *StringValueNodeSuite) TestSetCallsRegisteredSubscriberWithNewValue(c *check.C) {
	node := NewStringValueNode("", "init")
	var capturedValue string

	node.Subscribe(func(newValue string) {
		capturedValue = newValue
	})
	node.Set("new")

	c.Check(capturedValue, check.Equals, "new")
}
