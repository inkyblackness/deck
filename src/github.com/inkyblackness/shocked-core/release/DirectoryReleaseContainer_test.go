package release

import (
	"os"
	"path"
	"runtime"
	"sort"

	check "gopkg.in/check.v1"
)

type DirectoryReleaseContainerSuite struct {
	basePath string
	dynPath  string
}

var _ = check.Suite(&DirectoryReleaseContainerSuite{})

func (suite *DirectoryReleaseContainerSuite) SetUpSuite(c *check.C) {
	_, filename, _, _ := runtime.Caller(0)
	suite.basePath = path.Join(path.Dir(filename), "_test")
	suite.dynPath = path.Join(suite.basePath, "dynamicContainer")
	os.RemoveAll(suite.dynPath)
	os.MkdirAll(suite.dynPath, 0777)
}

func (suite *DirectoryReleaseContainerSuite) TearDownSuite(c *check.C) {
	if !c.Failed() {
		os.RemoveAll(suite.dynPath)
	}
}

func (suite *DirectoryReleaseContainerSuite) TestNewContainerFromDirReturnsErrorForNonExistingPath(c *check.C) {
	_, err := NewContainerFromDir(path.Join(suite.basePath, "not-existing-dir"))

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseContainerSuite) TestNamesReturnsDirectoryNames(c *check.C) {
	container, _ := NewContainerFromDir(path.Join(suite.basePath, "releases"))
	result := container.Names()
	sort.Strings(result)

	c.Check(result, check.DeepEquals, []string{"release1", "release2"})
}

func (suite *DirectoryReleaseContainerSuite) TestNamesReturnsOnlineDirectoryNames(c *check.C) {
	os.Mkdir(path.Join(suite.dynPath, "temp1"), 0755)
	container, _ := NewContainerFromDir(suite.dynPath)
	os.Mkdir(path.Join(suite.dynPath, "temp2"), 0755)
	result := container.Names()
	sort.Strings(result)

	c.Check(result, check.DeepEquals, []string{"temp1", "temp2"})
}

func (suite *DirectoryReleaseContainerSuite) TestGetReturnsErrorOnUnknownRelease(c *check.C) {
	container, _ := NewContainerFromDir(path.Join(suite.basePath, "releases"))
	_, err := container.Get("unknown-release")

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseContainerSuite) TestGetReturnsObjectForRelease(c *check.C) {
	container, _ := NewContainerFromDir(path.Join(suite.basePath, "releases"))
	rel, _ := container.Get("release2")

	c.Check(rel, check.NotNil)
}

func (suite *DirectoryReleaseContainerSuite) TestNewReturnsErrorForExistingRelease(c *check.C) {
	container, _ := NewContainerFromDir(path.Join(suite.basePath, "releases"))
	_, err := container.New("release1")

	c.Check(err, check.NotNil)
}

func (suite *DirectoryReleaseContainerSuite) TestNewReturnsObjectForNewRelease(c *check.C) {
	container, _ := NewContainerFromDir(path.Join(suite.basePath, "dynamicContainer"))
	rel, _ := container.New("releaseNew")
	info, _ := os.Stat(path.Join(suite.basePath, "dynamicContainer", "releaseNew"))

	c.Check(rel, check.NotNil)
	c.Check(info.IsDir(), check.Equals, true)
}
