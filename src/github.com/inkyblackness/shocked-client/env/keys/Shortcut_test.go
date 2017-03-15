package keys

import (
	check "gopkg.in/check.v1"
)

type ShortcutSuite struct{}

var _ = check.Suite(&ShortcutSuite{})

func (suite *ShortcutSuite) TestResolveShorcutReturnsValuesOfKnownCombo(c *check.C) {
	key, knownKey := ResolveShortcut("c", ModControl)

	c.Check(knownKey, check.Equals, true)
	c.Check(key, check.Equals, KeyCopy)
}

func (suite *ShortcutSuite) TestResolveShorcutReturnsFalseForUnknownCombo(c *check.C) {
	_, knownKey := ResolveShortcut("l", ModControl.With(ModShift))

	c.Check(knownKey, check.Equals, false)
}
