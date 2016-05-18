package release

import (
	"path"
	"runtime"

	check "gopkg.in/check.v1"
)

type AbsolutePathReleaseSuite struct {
	basePath  string
	goodPaths []string
}

var _ = check.Suite(&AbsolutePathReleaseSuite{})

func (suite *AbsolutePathReleaseSuite) SetUpSuite(c *check.C) {
	_, filename, _, _ := runtime.Caller(0)
	suite.basePath = path.Join(path.Dir(filename), "_test/absoluteBase")
	suite.goodPaths = []string{path.Join(suite.basePath, "sub1"), path.Join(suite.basePath, "sub2")}
}

func (suite *AbsolutePathReleaseSuite) TestFromAbsolutePathsReturnsErrorIfAPathDoesntExist(c *check.C) {
	_, err := FromAbsolutePaths(append(suite.goodPaths, path.Join(suite.basePath, "not-existing-dir")))

	c.Check(err, check.NotNil)
}

func (suite *AbsolutePathReleaseSuite) TestHasResourceReturnsTrueForExistingFiles(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)

	c.Check(rel.HasResource("file1.txt"), check.Equals, true)
	c.Check(rel.HasResource("file2.txt"), check.Equals, true)
}

func (suite *AbsolutePathReleaseSuite) TestHasResourceReturnsFalseForFilesInNestedDirectories(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)

	c.Check(rel.HasResource("file3.txt"), check.Equals, false)
}

func (suite *AbsolutePathReleaseSuite) TestHasResourceReturnsFalseForDirectories(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)

	c.Check(rel.HasResource("nested"), check.Equals, false)
}

func (suite *AbsolutePathReleaseSuite) TestGetResourceReturnsErrorForNonExistingFile(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)
	_, err := rel.GetResource("file5.txt")

	c.Check(err, check.NotNil)
}

func (suite *AbsolutePathReleaseSuite) TestGetResourceReturnsObjectForFile(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)
	resource, _ := rel.GetResource("file2.txt")

	c.Check(resource.Path(), check.Equals, path.Join(suite.basePath, "sub2"))
}

func (suite *AbsolutePathReleaseSuite) TestNewResourceReturnsErrorForExistingFile(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)
	_, err := rel.NewResource("file1.txt", path.Join(suite.basePath, "sub1"))

	c.Check(err, check.NotNil)
}

func (suite *AbsolutePathReleaseSuite) TestNewResourceReturnsErrorForUnknownPath(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)
	_, err := rel.NewResource("file10.txt", path.Join(suite.basePath, "unknown"))

	c.Check(err, check.NotNil)
}

func (suite *AbsolutePathReleaseSuite) TestNewResourceReturnsObjectForNewFile(c *check.C) {
	rel, _ := FromAbsolutePaths(suite.goodPaths)
	resource, _ := rel.NewResource("file10.txt", path.Join(suite.basePath, "sub1"))

	c.Check(resource.Name(), check.Equals, "file10.txt")
}
