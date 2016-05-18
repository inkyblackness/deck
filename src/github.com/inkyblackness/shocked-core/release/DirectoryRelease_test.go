package release

import (
	"os"
	"path"
	"runtime"

	check "gopkg.in/check.v1"
)

type DirectoryReleaseSuite struct {
	basePath string
	dynPath  string
}

var _ = check.Suite(&DirectoryReleaseSuite{})

func (suite *DirectoryReleaseSuite) SetUpSuite(c *check.C) {
	_, filename, _, _ := runtime.Caller(0)
	suite.basePath = path.Join(path.Dir(filename), "_test")
	suite.dynPath = path.Join(suite.basePath, "dynamicRelease")
	os.RemoveAll(suite.dynPath)
	os.MkdirAll(suite.dynPath, 0777)
}

func (suite *DirectoryReleaseSuite) TearDownSuite(c *check.C) {
	if !c.Failed() {
		os.RemoveAll(suite.dynPath)
	}
}

func (suite *DirectoryReleaseSuite) TestReleaseFromDirReturnsErrorIfDirectoryDoesntExist(c *check.C) {
	_, err := ReleaseFromDir(path.Join(suite.basePath, "not-existing-dir"))

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseSuite) TestReleaseFromDirReturnsErrorIfPathIsNotADir(c *check.C) {
	_, err := ReleaseFromDir(path.Join(suite.basePath, "testFile1.txt"))

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseSuite) TestHasResourceReturnsTrueForFileInSameDirectory(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))

	c.Check(rel.HasResource("file2.txt"), check.Equals, true)
}

func (suite *DirectoryReleaseSuite) TestHasResourceReturnsTrueForFileInNestedDirectory(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))

	c.Check(rel.HasResource("file4.txt"), check.Equals, true)
}

func (suite *DirectoryReleaseSuite) TestHasResourceReturnsFalseForFileInTooDeepNestedDirectory(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))

	c.Check(rel.HasResource("file5.txt"), check.Equals, false)
}

func (suite *DirectoryReleaseSuite) TestGetResourceReturnsErrorForNonExistingFile(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))
	_, err := rel.GetResource("file5.txt")

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseSuite) TestGetResourceReturnsObjectForFile(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))
	resource, _ := rel.GetResource("file3.txt")

	c.Check(resource.Path(), check.Equals, "nested1")
}

func (suite *DirectoryReleaseSuite) TestNewResourceReturnsErrorForExistingFile(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))
	_, err := rel.NewResource("file3.txt", ".")

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseSuite) TestNewResourceReturnsObjectForNewFile(c *check.C) {
	rel, _ := ReleaseFromDir(path.Join(suite.basePath, "releases", "release2"))
	resource, _ := rel.NewResource("file10.txt", "nested1")

	c.Check(resource.Name(), check.Equals, "file10.txt")
}
