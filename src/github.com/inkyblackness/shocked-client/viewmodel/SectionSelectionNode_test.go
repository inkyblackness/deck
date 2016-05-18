package viewmodel

import (
	"fmt"
	"math/rand"

	check "gopkg.in/check.v1"
)

type SectionSelectionNodeSuite struct {
}

var _ = check.Suite(&SectionSelectionNodeSuite{})

func (suite *SectionSelectionNodeSuite) someSectionNodes() map[string]*SectionNode {
	nodeMap := make(map[string]*SectionNode)
	count := rand.Intn(10)
	for index := 0; index < count; index++ {
		key := fmt.Sprintf("section%v", index)
		nodeMap[key] = NewSectionNode("Label-"+key, someNodeList(), NewBoolValueNode("l", true))
	}

	return nodeMap
}

func (suite *SectionSelectionNodeSuite) TestSpecializeCallsContainer(c *check.C) {
	node := NewSectionSelectionNode("", suite.someSectionNodes(), "")
	visitor := NewTestingNodeVisitor()

	node.Specialize(visitor)

	c.Check(visitor.sectionSelectionNodes, check.DeepEquals, []Node{node})
}

func (suite *SectionSelectionNodeSuite) TestLabel(c *check.C) {
	c.Check(NewSectionSelectionNode("l1", suite.someSectionNodes(), "").Label(), check.Equals, "l1")
	c.Check(NewSectionSelectionNode("l2", suite.someSectionNodes(), "").Label(), check.Equals, "l2")
}

func (suite *SectionSelectionNodeSuite) TestSections(c *check.C) {
	sections := suite.someSectionNodes()

	c.Check(NewSectionSelectionNode("l1", sections, "").Sections(), check.DeepEquals, sections)
}

func (suite *SectionSelectionNodeSuite) TestSelectionIsInitializedWithInitialValue(c *check.C) {
	sections := map[string]*SectionNode{
		"section1": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", true))}
	node := NewSectionSelectionNode("", sections, "section1")

	c.Check(node.Selection().Selected().Get(), check.Equals, "section1")
}

func (suite *SectionSelectionNodeSuite) TestSelectionContainsOnlyAvailableSections(c *check.C) {
	sections := map[string]*SectionNode{
		"section1": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", false)),
		"section2": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", true))}
	node := NewSectionSelectionNode("", sections, "section2")

	c.Check(node.Selection().Values(), check.DeepEquals, []string{"section2"})
}

func (suite *SectionSelectionNodeSuite) TestSelectionIsUpdatedWithChangingAvailabilityA(c *check.C) {
	sections := map[string]*SectionNode{
		"section1": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", false)),
		"section2": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", true))}
	node := NewSectionSelectionNode("", sections, "section2")

	sections["section1"].Available().Set(true)

	c.Check(len(node.Selection().Values()), check.Equals, 2)
}

func (suite *SectionSelectionNodeSuite) TestSelectionIsUpdatedWithChangingAvailabilityB(c *check.C) {
	sections := map[string]*SectionNode{
		"section1": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", false)),
		"section2": NewSectionNode("l1", someNodeList(), NewBoolValueNode("avail", true))}
	node := NewSectionSelectionNode("", sections, "section2")

	sections["section2"].Available().Set(false)

	c.Check(len(node.Selection().Values()), check.Equals, 0)
}
