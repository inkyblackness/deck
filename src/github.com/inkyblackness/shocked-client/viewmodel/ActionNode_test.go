package viewmodel

import (
	check "gopkg.in/check.v1"
)

type ActionNodeSuite struct {
}

var _ = check.Suite(&ActionNodeSuite{})

func (suite *ActionNodeSuite) TestSpecializeCallsAction(c *check.C) {
	node := NewActionNode("")
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.actionNodes, check.DeepEquals, []Node{node})
}

func (suite *ActionNodeSuite) TestLabel(c *check.C) {
	c.Check(NewActionNode("l1").Label(), check.Equals, "l1")
	c.Check(NewActionNode("l2").Label(), check.Equals, "l2")
}

func (suite *ActionNodeSuite) TestActCallsRegisteredSubscriber(c *check.C) {
	node := NewActionNode("")
	called := false

	node.Subscribe(func() {
		called = true
	})
	node.Act()

	c.Check(called, check.Equals, true)
}
