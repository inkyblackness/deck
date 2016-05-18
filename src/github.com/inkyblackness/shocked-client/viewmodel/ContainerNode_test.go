package viewmodel

import (
	check "gopkg.in/check.v1"
)

type ContainerNodeSuite struct {
}

var _ = check.Suite(&ContainerNodeSuite{})

func (suite *ContainerNodeSuite) TestSpecializeCallsContainer(c *check.C) {
	node := NewContainerNode("", someNodeMap())
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.containerNodes, check.DeepEquals, []Node{node})
}

func (suite *ContainerNodeSuite) TestLabel(c *check.C) {
	c.Check(NewContainerNode("l1", someNodeMap()).Label(), check.Equals, "l1")
	c.Check(NewContainerNode("l2", someNodeMap()).Label(), check.Equals, "l2")
}

func (suite *ContainerNodeSuite) TestGetReturnsInitialValue(c *check.C) {
	initial := someNodeMap()
	c.Check(NewContainerNode("", initial).Get(), check.DeepEquals, initial)
}
