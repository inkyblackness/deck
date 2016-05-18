package core

import (
	"fmt"
	"os"

	"github.com/inkyblackness/hacker/styling"

	check "gopkg.in/check.v1"
)

type HackerSuite struct {
	hacker *Hacker

	testDirectories map[string][]os.FileInfo
}

var _ = check.Suite(&HackerSuite{})

func (suite *HackerSuite) SetUpTest(c *check.C) {
	suite.testDirectories = make(map[string][]os.FileInfo)

	suite.hacker = NewHacker(styling.NullStyle())
	suite.hacker.fileAccess = fileAccess{
		readDir: func(path string) (info []os.FileInfo, err error) {
			var ok bool
			info, ok = suite.testDirectories[path]
			if !ok {
				err = fmt.Errorf("Not existing")
			}
			return
		}}

}

func (suite *HackerSuite) TestLoadOfUnknownLocationResultsInErrorMessage(c *check.C) {
	result := suite.hacker.Load("nonExisting1", "nonExisting2")

	c.Check(result, check.Equals, "Can't access directories")
}

func (suite *HackerSuite) TestLoadOfWrongLocationResultsInErrorMessage(c *check.C) {
	suite.testDirectories["dir1"] = []os.FileInfo{testFile("file1.res"), testFile("file2.res")}
	suite.testDirectories["dir2"] = []os.FileInfo{testFile("file3.res"), testFile("file4.res")}

	result := suite.hacker.Load("dir1", "dir2")

	c.Check(result, check.Equals, "Could not resolve release")
}

func (suite *HackerSuite) TestLoadOfKnownLocationResultsInConfirmation(c *check.C) {
	hdFiles, cdFiles := DataFiles(&dosCdRelease)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.testDirectories["dir2"] = testFiles(cdFiles...)

	result := suite.hacker.Load("dir1", "dir2")

	c.Check(result, check.Equals, "Loaded release [DOS CD Release]")
}

func (suite *HackerSuite) TestLoadAllowsOptionalSecondPath(c *check.C) {
	hdFiles, _ := DataFiles(&dosHdDemo)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)

	result := suite.hacker.Load("dir1", "")

	c.Check(result, check.Equals, "Loaded release [DOS HD Demo]")
}

func (suite *HackerSuite) TestLoadOfKnownSwitchedLocationResultsInConfirmation(c *check.C) {
	hdFiles, cdFiles := DataFiles(&dosCdDemo)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.testDirectories["dir2"] = testFiles(cdFiles...)

	result := suite.hacker.Load("dir2", "dir1")

	c.Check(result, check.Equals, "Loaded release [DOS CD Demo]")
}

func (suite *HackerSuite) TestLoadSetsUpRootNodeForHdOnly(c *check.C) {
	hdFiles, _ := DataFiles(&dosHdDemo)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.hacker.Load("dir1", "")

	c.Assert(suite.hacker.root, check.Not(check.IsNil))
	suite.checkLocationHasDir(c, HD, "dir1")
	//c.Check(suite.hacker.root.Resolve(HD.String()).filePath, check.Equals, "dir1")
	c.Check(len(suite.hacker.root.Children()), check.Equals, 1)
}

func (suite *HackerSuite) checkLocationHasDir(c *check.C, location DataLocation, expected string) {
	node := suite.hacker.root.Resolve(location.String()).(*locationDataNode)

	c.Assert(node, check.Not(check.IsNil))
	c.Check(node.filePath, check.Equals, expected)
}

func (suite *HackerSuite) TestLoadSetsUpRootNodeForRelease(c *check.C) {
	hdFiles, cdFiles := DataFiles(&dosCdRelease)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.testDirectories["dir2"] = testFiles(cdFiles...)
	suite.hacker.Load("dir1", "dir2")

	c.Assert(suite.hacker.root, check.Not(check.IsNil))
	suite.checkLocationHasDir(c, HD, "dir1")
	suite.checkLocationHasDir(c, CD, "dir2")
}

func (suite *HackerSuite) TestLoadSetsUpRootNodeForSwappedPaths(c *check.C) {
	hdFiles, cdFiles := DataFiles(&dosCdRelease)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.testDirectories["dir2"] = testFiles(cdFiles...)
	suite.hacker.Load("dir2", "dir1")

	c.Assert(suite.hacker.root, check.Not(check.IsNil))
	suite.checkLocationHasDir(c, HD, "dir1")
	suite.checkLocationHasDir(c, CD, "dir2")
}

func (suite *HackerSuite) TestInfoWithoutDataReturnsHintToLoad(c *check.C) {
	result := suite.hacker.Info()

	c.Check(result, check.Equals, `No data loaded. Use the [load "path1" "path2"] command.`)
}

func (suite *HackerSuite) givenAStandardSetup() {
	hdFiles, cdFiles := DataFiles(&dosCdRelease)
	suite.testDirectories["dir1"] = testFiles(hdFiles...)
	suite.testDirectories["dir2"] = testFiles(cdFiles...)
	suite.hacker.Load("dir1", "dir2")
}

func (suite *HackerSuite) TestInfoAfterLoadReturnsReleaseInfo(c *check.C) {
	suite.givenAStandardSetup()

	result := suite.hacker.Info()

	c.Check(result, check.Equals, suite.hacker.root.Info())
}

func (suite *HackerSuite) TestChangeDirectoryChangesCurrentNode(c *check.C) {
	suite.givenAStandardSetup()

	suite.hacker.ChangeDirectory("hd")

	c.Check(suite.hacker.Info(), check.Equals, suite.hacker.root.Resolve(HD.String()).Info())
}

func (suite *HackerSuite) TestChangeDirectoryHandlesStartingSlash(c *check.C) {
	suite.givenAStandardSetup()
	suite.hacker.ChangeDirectory("hd")

	suite.hacker.ChangeDirectory("/cd")

	c.Check(suite.hacker.Info(), check.Equals, suite.hacker.root.Resolve(CD.String()).Info())
}

func (suite *HackerSuite) TestChangeDirectoryHandlesDotDot(c *check.C) {
	suite.givenAStandardSetup()
	suite.hacker.ChangeDirectory("hd")

	suite.hacker.ChangeDirectory("../cd")

	c.Check(suite.hacker.Info(), check.Equals, suite.hacker.root.Resolve(CD.String()).Info())
}

func (suite *HackerSuite) TestChangeDirectoryIgnoresTrailingSlash(c *check.C) {
	suite.givenAStandardSetup()

	suite.hacker.ChangeDirectory("hd/")

	c.Check(suite.hacker.Info(), check.Equals, suite.hacker.root.Resolve(HD.String()).Info())
}

func (suite *HackerSuite) TestCurrentDirctoryReturnsCurrentPath(c *check.C) {
	suite.givenAStandardSetup()

	suite.hacker.ChangeDirectory("hd")

	c.Check(suite.hacker.CurrentDirectory(), check.Equals, "/hd")
}

func (suite *HackerSuite) TestDumpReturnsCurrentDataInDumpFormat(c *check.C) {
	dataNode := NewTestingDataNode("test")
	suite.givenAStandardSetup()
	suite.hacker.curNode = dataNode
	dataNode.data = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x41}

	result := suite.hacker.Dump()

	c.Check(result, check.Equals, "0000  00 01 02 03 04 05 06 07  08 09 0A 0B 0C 0D 0E 0F  ........ ........\n"+
		"0010  41                                                A                \n")
}

func (suite *HackerSuite) TestDiffOfNodesWithoutDataComparesChildren(c *check.C) {
	parent1 := NewTestingDataNode("parent1")
	parent1.addChild(NewTestingDataNode("child1"))
	parent2 := NewTestingDataNode("parent2")
	suite.givenAStandardSetup()

	suite.hacker.root.addChild(parent1)
	suite.hacker.curNode = parent2

	result := suite.hacker.Diff("/parent1")

	c.Check(result, check.Equals, "- /parent1/child1\n")
}

func (suite *HackerSuite) TestDiffNodesReportsRemovedChildren(c *check.C) {
	parent1 := NewTestingDataNode("parent1")
	parent1.addChild(NewTestingDataNode("child1"))
	parent2 := NewTestingDataNode("parent2")

	result := suite.hacker.diffNodes("/parent1", parent1, "/parent2", parent2)

	c.Check(result, check.Equals, "- /parent1/child1\n")
}

func (suite *HackerSuite) TestDiffNodesReportsAddedChildren(c *check.C) {
	parent1 := NewTestingDataNode("parent1")
	parent2 := NewTestingDataNode("parent2")
	parent2.addChild(NewTestingDataNode("child1"))

	result := suite.hacker.diffNodes("/parent1", parent1, "/parent2", parent2)

	c.Check(result, check.Equals, "+ /parent2/child1\n")
}

func (suite *HackerSuite) TestDiffNodesReportsChangeInData(c *check.C) {
	parent1 := NewTestingDataNode("parent1")
	child11 := NewTestingDataNode("child1")
	child11.data = []byte{0x01, 0x02}
	parent1.addChild(child11)
	parent2 := NewTestingDataNode("parent2")
	child21 := NewTestingDataNode("child1")
	child21.data = []byte{0x02, 0x03}
	parent2.addChild(child21)

	result := suite.hacker.diffNodes("/parent1", parent1, "/parent2", parent2)

	c.Check(result, check.Equals, "M /parent2/child1\n")
}

func (suite *HackerSuite) TestDiffNodesReportsChangeInDataRecursive(c *check.C) {
	parent1 := NewTestingDataNode("parent1")
	child11 := NewTestingDataNode("child1")
	child111 := NewTestingDataNode("child1.1")
	child111.data = []byte{0x01, 0x02}
	child11.addChild(child111)
	parent1.addChild(child11)
	parent2 := NewTestingDataNode("parent2")
	child21 := NewTestingDataNode("child1")
	child211 := NewTestingDataNode("child1.1")
	child211.data = []byte{0x02, 0x03}
	child21.addChild(child211)
	parent2.addChild(child21)

	result := suite.hacker.diffNodes("/parent1", parent1, "/parent2", parent2)

	c.Check(result, check.Equals, "M /parent2/child1/child1.1\n")
}

func (suite *HackerSuite) TestPutWorksWithDataNodes(c *check.C) {

	node := NewTestingDataNode("rawNode")
	node.data = []byte{0x01, 0x02, 0x03, 0x04}

	suite.hacker.curNode = node
	suite.hacker.Put(0, []byte{0x0A, 0x0B})

	c.Check(node.data, check.DeepEquals, []byte{0x0A, 0x0B, 0x03, 0x04})
}
