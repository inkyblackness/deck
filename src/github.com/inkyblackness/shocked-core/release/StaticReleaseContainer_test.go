package release

import (
	"sort"

	check "gopkg.in/check.v1"
)

type StaticReleaseContainerSuite struct {
	basePath string
	dynPath  string
}

var _ = check.Suite(&StaticReleaseContainerSuite{})

func (suite *StaticReleaseContainerSuite) SetUpSuite(c *check.C) {

}

func (suite *StaticReleaseContainerSuite) TestNewReturnsError(c *check.C) {
	container := NewStaticReleaseContainer(map[string]Release{})
	_, err := container.New("test1")

	c.Check(err, check.NotNil)
}

func (suite *StaticReleaseContainerSuite) TestNamesReturnsListOfRegisteredReleases(c *check.C) {
	rel1 := NewMemoryRelease()
	rel2 := NewMemoryRelease()
	container := NewStaticReleaseContainer(map[string]Release{"rel1": rel1, "rel2": rel2})

	names := container.Names()
	sort.Strings(names)

	c.Check(names, check.DeepEquals, []string{"rel1", "rel2"})
}

func (suite *StaticReleaseContainerSuite) TestGetReturnsRegisteredReleases(c *check.C) {
	rel1 := NewMemoryRelease()
	rel2 := NewMemoryRelease()
	container := NewStaticReleaseContainer(map[string]Release{"rel1": rel1, "rel2": rel2})

	retrieved1, _ := container.Get("rel1")
	retrieved2, _ := container.Get("rel2")

	c.Check(retrieved1, check.Equals, rel1)
	c.Check(retrieved2, check.Equals, rel2)
}

func (suite *StaticReleaseContainerSuite) TestGetReturnsErrorForUnknownRelease(c *check.C) {
	container := NewStaticReleaseContainer(map[string]Release{})

	_, err := container.Get("rel3")

	c.Check(err, check.NotNil)
}
