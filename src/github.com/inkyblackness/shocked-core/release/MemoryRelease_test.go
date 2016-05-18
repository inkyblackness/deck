package release

import (
	check "gopkg.in/check.v1"
)

type MemoryReleaseSuite struct {
	release Release
}

var _ = check.Suite(&MemoryReleaseSuite{})

func (suite *MemoryReleaseSuite) SetUpTest(c *check.C) {
	suite.release = NewMemoryRelease()
}

func (suite *MemoryReleaseSuite) TestNewResourceReturnsResourceIfNotExisting(c *check.C) {
	resource, err := suite.release.NewResource("test1.res", "rel")

	c.Assert(err, check.IsNil)
	c.Check(resource, check.NotNil)
}

func (suite *MemoryReleaseSuite) TestNewResourceReturnsErrorIfExisting(c *check.C) {
	suite.release.NewResource("test1.res", "rel1")
	resource, err := suite.release.NewResource("test1.res", "rel2")

	c.Assert(resource, check.IsNil)
	c.Check(err, check.NotNil)
}

func (suite *MemoryReleaseSuite) TestHasResourceReturnsFalseIfNotExisting(c *check.C) {
	result := suite.release.HasResource("test1.res")

	c.Check(result, check.Equals, false)
}

func (suite *MemoryReleaseSuite) TestHasResourceReturnsTrueIfExisting(c *check.C) {
	suite.release.NewResource("other.res", "rel1")
	result := suite.release.HasResource("other.res")

	c.Check(result, check.Equals, true)
}

func (suite *MemoryReleaseSuite) TestGetResourceReturnsErrorIfNotExisting(c *check.C) {
	resource, err := suite.release.GetResource("test1.res")

	c.Assert(resource, check.IsNil)
	c.Check(err, check.NotNil)
}

func (suite *MemoryReleaseSuite) TestGetResourceReturnsResourceExisting(c *check.C) {
	suite.release.NewResource("test1.res", "rel1")
	resource, err := suite.release.GetResource("test1.res")

	c.Assert(err, check.IsNil)
	c.Check(resource, check.NotNil)
}
