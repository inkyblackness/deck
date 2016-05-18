package viewmodel

import (
	check "gopkg.in/check.v1"
)

type TableNodeSuite struct {
}

var _ = check.Suite(&TableNodeSuite{})

func (suite *TableNodeSuite) TestLabel(c *check.C) {
	c.Check(NewTableNode("l1").Label(), check.Equals, "l1")
	c.Check(NewTableNode("l2").Label(), check.Equals, "l2")
}

func (suite *TableNodeSuite) TestSpecializeCallsTable(c *check.C) {
	node := NewTableNode("label")
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.tableNodes, check.DeepEquals, []Node{node})
}

func (suite *TableNodeSuite) TestGetReturnsInitialValue(c *check.C) {
	initial := NewContainerNode("label", someNodeMap())
	c.Check(NewTableNode("label", initial).Get(), check.DeepEquals, []*ContainerNode{initial})
}

func (suite *TableNodeSuite) TestSetChangesCurrentValue(c *check.C) {
	node := NewTableNode("label")
	rows := []*ContainerNode{NewContainerNode("l1", someNodeMap())}

	node.Set(rows)

	c.Check(node.Get(), check.DeepEquals, rows)
}

func (suite *TableNodeSuite) TestSetCallsRegisteredSubscriberWithNewEntries(c *check.C) {
	node := NewTableNode("")
	var capturedRows []*ContainerNode

	node.Subscribe(func(newRows []*ContainerNode) {
		capturedRows = newRows
	})

	newEntry := NewContainerNode("l1", someNodeMap())
	node.Set([]*ContainerNode{newEntry})

	c.Check(capturedRows, check.DeepEquals, []*ContainerNode{newEntry})
}
