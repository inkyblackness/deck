package core

import (
	check "gopkg.in/check.v1"
)

type RootDataNodeSuite struct {
	node *rootDataNode
}

var _ = check.Suite(&RootDataNodeSuite{})

func (suite *RootDataNodeSuite) SetUpTest(c *check.C) {
	suite.node = newRootDataNode(&dosCdRelease)
	suite.node.addChild(newLocationDataNode(suite.node, HD, "/hdPath/", []string{"file1", "file2"}, nil))
	suite.node.addChild(newLocationDataNode(suite.node, CD, "/cdPath/", []string{"fileA", "fileB"}, nil))
}

func (suite *RootDataNodeSuite) TestResolveOfDotDotReturnsNil(c *check.C) {
	result := suite.node.Resolve("..")

	c.Check(result, check.IsNil)
}

func (suite *RootDataNodeSuite) TestResolveOfHdReturnsHdLocation(c *check.C) {
	result := suite.node.Resolve("hd")

	c.Check(result.ID(), check.Equals, "hd")
}

func (suite *RootDataNodeSuite) TestResolveOfUnknownEntryReturnsNil(c *check.C) {
	result := suite.node.Resolve("hw")

	c.Check(result, check.IsNil)
}
