package interpreters

import (
	check "gopkg.in/check.v1"
)

type DescriptionSuite struct {
}

var _ = check.Suite(&DescriptionSuite{})

func (suite *DescriptionSuite) TestWithReturnsANewDescription(c *check.C) {
	second := New().With("test", 0, 32)

	c.Assert(second, check.NotNil)
	c.Check(second, check.Not(check.DeepEquals), New())
}

func (suite *DescriptionSuite) TestWithCopiesPreviousFields(c *check.C) {
	first := New().With("field1", 0, 32)
	second := first.With("field2", 32, 32)
	onlyField2 := New().With("field2", 32, 32)

	c.Check(second, check.Not(check.DeepEquals), onlyField2)
}

func (suite *DescriptionSuite) TestWithLeavesOriginalAlone(c *check.C) {
	first := New().With("field1", 0, 32)
	onlyField1 := New().With("field1", 0, 32)

	first.With("field2", 32, 32)

	c.Check(first, check.DeepEquals, onlyField1)
}

func (suite *DescriptionSuite) TestForCreatesNewInstance(c *check.C) {
	data := make([]byte, 10)
	inst := New().For(data)

	c.Check(inst, check.NotNil)
}

func (suite *DescriptionSuite) TestRefiningReturnsANewDescription(c *check.C) {
	refined := New().With("field1", 0, 16)
	second := New().Refining("test", 0, 4, refined, Always)

	c.Assert(second, check.NotNil)
	c.Check(second, check.Not(check.DeepEquals), New())
}

func (suite *DescriptionSuite) TestRefiningCopiesPreviousFields(c *check.C) {
	first := New().With("fieldA", 0, 8)
	refined := New().With("field1", 0, 16)
	second := first.Refining("test", 0, 4, refined, Always)
	secondMissing := New().Refining("test", 0, 4, refined, Always)

	c.Check(second, check.Not(check.DeepEquals), secondMissing)
}

func (suite *DescriptionSuite) TestRefiningCopiesPreviousRefinements(c *check.C) {
	first := New().Refining("sub1", 0, 1, New(), Always)
	second := first.Refining("sub2", 0, 1, New(), Always)

	c.Check(second.refinements["sub1"], check.NotNil)
}

func (suite *DescriptionSuite) TestAsPanicsIfNoFieldIsActive(c *check.C) {
	c.Check(func() { New().As(func(simpl *Simplifier) bool { return false }) }, check.PanicMatches, "No field active")
}

func (suite *DescriptionSuite) TestAsRegistersRangeFunctionForLastField(c *check.C) {
	rangeFuncCalled := false
	rangeFunc := func(simpl *Simplifier) bool {
		rangeFuncCalled = true
		return true
	}
	desc := New().With("fieldA", 0, 2).As(rangeFunc)

	inst := desc.For([]byte{0x00, 0x01})
	simplifier := NewSimplifier(func(minValue, maxValue int64, formatter RawValueFormatter) {})
	inst.Describe("fieldA", simplifier)

	c.Check(rangeFuncCalled, check.Equals, true)
}
