package ui

import (
	check "gopkg.in/check.v1"
)

type AreaBuilderSuite struct {
	builder *AreaBuilder
}

var _ = check.Suite(&AreaBuilderSuite{})

func (suite *AreaBuilderSuite) SetUpTest(c *check.C) {
	suite.builder = NewAreaBuilder()
}

func (suite *AreaBuilderSuite) TestBuildCreatesAnArea(c *check.C) {
	area := suite.builder.Build()

	c.Check(area, check.NotNil)
}

func (suite *AreaBuilderSuite) TestPositionAnchorsDefaultToZero(c *check.C) {
	area := suite.builder.Build()

	c.Check(area.Left().Value(), check.Equals, float32(0.0))
	c.Check(area.Top().Value(), check.Equals, float32(0.0))
	c.Check(area.Right().Value(), check.Equals, float32(0.0))
	c.Check(area.Bottom().Value(), check.Equals, float32(0.0))
}

func (suite *AreaBuilderSuite) TestPositionAnchorsAreaTakenOver(c *check.C) {
	suite.builder.SetLeft(NewAbsoluteAnchor(10.0))
	suite.builder.SetTop(NewAbsoluteAnchor(20.0))
	suite.builder.SetRight(NewAbsoluteAnchor(30.0))
	suite.builder.SetBottom(NewAbsoluteAnchor(40.0))
	area := suite.builder.Build()

	c.Check(area.Left().Value(), check.Equals, float32(10.0))
	c.Check(area.Top().Value(), check.Equals, float32(20.0))
	c.Check(area.Right().Value(), check.Equals, float32(30.0))
	c.Check(area.Bottom().Value(), check.Equals, float32(40.0))
}

func (suite *AreaBuilderSuite) TestDefaultRenderFunctionIsSet(c *check.C) {
	area := suite.builder.Build()

	c.Check(area.onRender, check.NotNil)
}

func (suite *AreaBuilderSuite) TestOnRenderSetsRenderFunction(c *check.C) {
	called := false
	onRender := func(*Area) { called = true }
	suite.builder.OnRender(onRender)
	area := suite.builder.Build()

	area.Render()

	c.Check(called, check.Equals, true)
}
