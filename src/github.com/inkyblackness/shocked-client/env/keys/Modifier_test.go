package keys

import (
	check "gopkg.in/check.v1"
)

type ModifierSuite struct{}

var _ = check.Suite(&ModifierSuite{})

func (suite *ModifierSuite) TestHasReturnsTrueForThemselves(c *check.C) {
	mods := []Modifier{ModShift, ModControl, ModAlt, ModSuper}

	for _, mod := range mods {
		c.Check(mod.Has(mod), check.Equals, true)
	}
}

func (suite *ModifierSuite) TestHasReturnsFalseForDifferentModifiers(c *check.C) {
	mods := []Modifier{ModShift, ModControl, ModAlt, ModSuper}

	for index, mod := range mods {
		nextMod := mods[(index+1)%len(mods)]
		c.Check(mod.Has(nextMod), check.Equals, false)
	}
}

func (suite *ModifierSuite) TestHasReturnsFalseForSetsIncludingOthers(c *check.C) {
	mods := []Modifier{ModShift, ModControl, ModAlt, ModSuper}

	for index, mod := range mods {
		nextMod := mods[(index+1)%len(mods)]
		c.Check(mod.Has(mod.With(nextMod)), check.Equals, false)
	}
}

func (suite *ModifierSuite) TestWithoutReturnsReduction(c *check.C) {
	mod := ModShift.With(ModAlt).Without(ModShift)

	c.Check(mod, check.Equals, ModAlt)
}
