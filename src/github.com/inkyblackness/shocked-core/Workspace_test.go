package core

import (
	"sort"

	"github.com/inkyblackness/shocked-core/release"

	check "gopkg.in/check.v1"
)

type WorkspaceSuite struct {
	source   release.Release
	projects release.ReleaseContainer
}

var _ = check.Suite(&WorkspaceSuite{})

func (suite *WorkspaceSuite) SetUpTest(c *check.C) {
	suite.source = release.NewMemoryRelease()
	suite.projects = release.NewMemoryReleaseContainer()
	suite.projects.New("project1")
	suite.projects.New("project2")
}

func (suite *WorkspaceSuite) TestNamesReturnsNamesOfProjects(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	result := ws.ProjectNames()
	sort.Strings(result)

	c.Check(result, check.DeepEquals, []string{"project1", "project2"})
}

func (suite *WorkspaceSuite) TestProjectReturnsErrorIfNoReleaseAvailable(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	_, err := ws.Project("not-existing")

	c.Check(err, check.NotNil)
}

func (suite *WorkspaceSuite) TestProjectReturnsInstanceIfReleaseAvailable(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	project, _ := ws.Project("project2")

	c.Check(project, check.NotNil)
}

func (suite *WorkspaceSuite) TestNewProjectReturnsErrorIfExisting(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	_, err := ws.NewProject("project1")

	c.Check(err, check.NotNil)
}

func (suite *WorkspaceSuite) TestNewProjectReturnsInstanceIfNotExisting(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	project, _ := ws.NewProject("project3")

	c.Check(project, check.NotNil)
}

func (suite *WorkspaceSuite) TestNewProjectCreatesNewRelease(c *check.C) {
	ws := NewWorkspace(suite.source, suite.projects)

	ws.NewProject("project3")
	rel, _ := suite.projects.Get("project3")

	c.Check(rel, check.NotNil)
}
