package release

import (
	"os"
	"path"
	"runtime"

	check "gopkg.in/check.v1"
)

type FileResourceSuite struct {
	basePath string
	dynPath  string
}

var _ = check.Suite(&FileResourceSuite{})

func (suite *FileResourceSuite) SetUpSuite(c *check.C) {
	_, filename, _, _ := runtime.Caller(0)
	suite.basePath = path.Join(path.Dir(filename), "_test")
	suite.dynPath = path.Join(suite.basePath, "dynamic")
	os.RemoveAll(suite.dynPath)
	os.MkdirAll(suite.dynPath, 0777)
}

func (suite *FileResourceSuite) TearDownSuite(c *check.C) {
	if !c.Failed() {
		os.RemoveAll(suite.dynPath)
	}
}

func (suite *FileResourceSuite) TestNameReturnsNameOfResource(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, "rel", "file")

	c.Check(resource.Name(), check.Equals, "test1.res")
}

func (suite *FileResourceSuite) TestPathReturnsRelativePathOfResource(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, "rel", "file")

	c.Check(resource.Path(), check.Equals, "rel")
}

func (suite *FileResourceSuite) TestAsSourceReturnsErrorForNotExisting(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, ".", "notExisting.bin")
	_, err := resource.AsSource()

	c.Check(err, check.NotNil)
}

func (suite *FileResourceSuite) TestAsSourceReturnsFileIfExisting(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, ".", "testFile1.txt")
	file, err := resource.AsSource()
	c.Assert(err, check.IsNil)
	defer file.Close()

	c.Check(file, check.NotNil)
}

func (suite *FileResourceSuite) TestAsSinkReturnsErrorForNotWritable(c *check.C) {
	lockedFile, lockedErr := os.OpenFile(path.Join(suite.basePath, "dynamic", "locked.bin"), os.O_CREATE|os.O_RDWR|0777, os.ModeExclusive)
	c.Assert(lockedErr, check.IsNil)
	defer lockedFile.Close()

	resource := newFileResource("test1.res", suite.basePath, "dynamic", "locked.bin")
	_, err := resource.AsSink()

	c.Check(err, check.NotNil)
}

func (suite *FileResourceSuite) TestAsSinkReturnsObjectForNewFile(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, "dynamic", "newFile.bin")
	file, err := resource.AsSink()
	c.Assert(err, check.IsNil)
	defer file.Close()

	c.Check(file, check.NotNil)
}

func (suite *FileResourceSuite) TestAsSinkCreatesDirectoriesForNewFile(c *check.C) {
	resource := newFileResource("test1.res", suite.basePath, "dynamic/created", "newFile.bin")
	file, err := resource.AsSink()
	c.Assert(err, check.IsNil)
	defer file.Close()

	c.Check(file, check.NotNil)
}
