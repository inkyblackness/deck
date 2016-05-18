package core

import (
	check "gopkg.in/check.v1"
)

type ReleaseDescSuite struct {
}

var _ = check.Suite(&ReleaseDescSuite{})

func (suite *ReleaseDescSuite) SetUpTest(c *check.C) {

}

func (suite *ReleaseDescSuite) TestFindReleaseReturnsNilForUnknown(c *check.C) {
	result := FindRelease(nil, nil)

	c.Check(result, check.IsNil)
}

func (suite *ReleaseDescSuite) TestFindReleaseCanDetermineAllReleases(c *check.C) {
	for _, release := range Releases {
		hdFiles, cdFiles := DataFiles(release)

		result := FindRelease(hdFiles, cdFiles)

		c.Check(result, check.Equals, release)
	}
}
