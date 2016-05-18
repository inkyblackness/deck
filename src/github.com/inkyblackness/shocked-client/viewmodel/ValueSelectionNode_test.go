package viewmodel

import (
	check "gopkg.in/check.v1"
)

type ValueSelectionNodeSuite struct {
}

var _ = check.Suite(&ValueSelectionNodeSuite{})

func (suite *ValueSelectionNodeSuite) TestSpecializeCallsValueSelection(c *check.C) {
	node := NewValueSelectionNode("", []string{}, "")
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.valueSelectionNodes, check.DeepEquals, []Node{node})
}

func (suite *ValueSelectionNodeSuite) TestLabel(c *check.C) {
	c.Check(NewValueSelectionNode("l1", []string{}, "").Label(), check.Equals, "l1")
	c.Check(NewValueSelectionNode("l2", []string{}, "").Label(), check.Equals, "l2")
}

func (suite *ValueSelectionNodeSuite) TestSelectedIsInitialized(c *check.C) {
	node := NewValueSelectionNode("", []string{"a", "b"}, "a")

	c.Check(node.Selected().Get(), check.Equals, "a")
}

func (suite *ValueSelectionNodeSuite) TestValuesReturnsInitialValue(c *check.C) {
	values := []string{"a", "b"}
	node := NewValueSelectionNode("", values, "")

	c.Check(node.Values(), check.DeepEquals, values)
}

func (suite *ValueSelectionNodeSuite) TestSetValuesModifiesValues(c *check.C) {
	node := NewValueSelectionNode("", []string{"a"}, "")

	newValues := []string{"b", "c"}
	node.SetValues(newValues)

	c.Check(node.Values(), check.DeepEquals, newValues)
}

func (suite *ValueSelectionNodeSuite) TestSetValuesNotifiesListener(c *check.C) {
	node := NewValueSelectionNode("", []string{"a"}, "")
	var capturedValues []string

	node.Subscribe(func(newValues []string) {
		capturedValues = newValues
	})
	newValues := []string{"b", "c"}
	node.SetValues(newValues)

	c.Check(capturedValues, check.DeepEquals, newValues)
}

func (suite *ValueSelectionNodeSuite) TestSetValuesResetsSelectedIfValueLost(c *check.C) {
	node := NewValueSelectionNode("", []string{"a"}, "a")

	node.SetValues([]string{"b"})

	c.Check(node.Selected().Get(), check.Equals, "")
}

func (suite *ValueSelectionNodeSuite) TestSetValuesKeepsSelectedIfValueStillThere(c *check.C) {
	node := NewValueSelectionNode("", []string{"a"}, "a")

	node.SetValues([]string{"b", "c", "a"})

	c.Check(node.Selected().Get(), check.Equals, "a")
}
