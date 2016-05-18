package release

import (
	"sort"

	check "gopkg.in/check.v1"
)

type MemoryReleaseContainerSuite struct {
	container ReleaseContainer
}

var _ = check.Suite(&MemoryReleaseContainerSuite{})

func (suite *MemoryReleaseContainerSuite) SetUpTest(c *check.C) {
	suite.container = NewMemoryReleaseContainer()
}

func (suite *MemoryReleaseContainerSuite) TestNewReturnsNewReleaseIfNotExisting(c *check.C) {
	rel, err := suite.container.New("rel1")

	c.Assert(err, check.IsNil)
	c.Check(rel, check.NotNil)
}

func (suite *MemoryReleaseContainerSuite) TestNewReturnsErrorIfExisting(c *check.C) {
	suite.container.New("rel2")
	rel, err := suite.container.New("rel2")

	c.Assert(rel, check.IsNil)
	c.Check(err, check.NotNil)
}

func (suite *MemoryReleaseContainerSuite) TestGetReturnsErrorIfNotExisting(c *check.C) {
	rel, err := suite.container.Get("rel1")

	c.Assert(rel, check.IsNil)
	c.Check(err, check.NotNil)
}

func (suite *MemoryReleaseContainerSuite) TestGetReturnsReleaseIfExisting(c *check.C) {
	suite.container.New("rel3")
	rel, err := suite.container.Get("rel3")

	c.Assert(err, check.IsNil)
	c.Check(rel, check.NotNil)
}

func (suite *MemoryReleaseContainerSuite) TestNamesReturnsEmptyArrayIfEmpty(c *check.C) {
	names := suite.container.Names()

	c.Check(len(names), check.Equals, 0)
}

func (suite *MemoryReleaseContainerSuite) TestNamesReturnsNamesIfNotEmpty(c *check.C) {
	suite.container.New("rel3")
	suite.container.New("rel1")
	names := suite.container.Names()
	sort.Strings(names)

	c.Check(names, check.DeepEquals, []string{"rel1", "rel3"})
}
