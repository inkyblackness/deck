package cmd

import (
	check "gopkg.in/check.v1"
)

type CombinedSourceSuite struct {
}

var _ = check.Suite(&CombinedSourceSuite{})

func (suite *CombinedSourceSuite) SetUpTest(c *check.C) {
}

func (suite *CombinedSourceSuite) TestNextOnEmptySourceImmediatelyFinishes(c *check.C) {
	combined := NewCombinedSource()

	_, finished := combined.Next()

	c.Assert(finished, check.Equals, true)
}

func (suite *CombinedSourceSuite) TestNextOnSingleSourceReturnsSourceData(c *check.C) {
	source := NewStaticSource("test")
	combined := NewCombinedSource(source)

	cmd, _ := combined.Next()

	c.Assert(cmd, check.Equals, "test")
}

func (suite *CombinedSourceSuite) TestNextSkipsToSecondSourceIfFirstFinished(c *check.C) {
	source := NewStaticSource("second")
	combined := NewCombinedSource(NewStaticSource(), source)

	cmd, _ := combined.Next()

	c.Assert(cmd, check.Equals, "second")
}

func (suite *CombinedSourceSuite) TestNextSkipsToThirdSourceIfPreviousFinished(c *check.C) {
	source := NewStaticSource("third")
	combined := NewCombinedSource(NewStaticSource(), NewStaticSource(), source)

	cmd, _ := combined.Next()

	c.Assert(cmd, check.Equals, "third")
}

func (suite *CombinedSourceSuite) TestNextFinishesWhenLastSourceFinished(c *check.C) {
	combined := NewCombinedSource(NewStaticSource())

	_, finished := combined.Next()

	c.Assert(finished, check.Equals, true)
}
