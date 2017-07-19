package interpreters

import (
	check "gopkg.in/check.v1"
)

type SimplifierSuite struct {
}

var _ = check.Suite(&SimplifierSuite{})

func (suite *SimplifierSuite) SetUpTest(c *check.C) {
}

func (suite *SimplifierSuite) TestRawValueCallsHandler(c *check.C) {
	var calledMinValue int64
	var calledMaxValue int64
	rawHandler := func(minValue, maxValue int64, formatter RawValueFormatter) {
		calledMinValue, calledMaxValue = minValue, maxValue
	}
	simpl := NewSimplifier(rawHandler)

	simpl.rawValue(&entry{count: 2})

	c.Check(calledMinValue, check.Equals, int64(-1))
	c.Check(calledMaxValue, check.Equals, int64(32767))
}

func (suite *SimplifierSuite) TestEnumValueReturnsFalseIfNoHandlerRegistered(c *check.C) {
	simpl := NewSimplifier(func(minValue, maxValue int64, formatter RawValueFormatter) {})

	result := simpl.enumValue(map[uint32]string{})

	c.Check(result, check.Equals, false)
}

func (suite *SimplifierSuite) TestEnumValueCallsRegisteredHandler(c *check.C) {
	result := map[uint32]string{}
	simpl := NewSimplifier(func(minValue, maxValue int64, formatter RawValueFormatter) {})
	simpl.SetEnumValueHandler(func(values map[uint32]string) {
		result = values
	})

	expected := map[uint32]string{1: "value"}
	simpl.enumValue(expected)

	c.Check(result, check.DeepEquals, expected)
}
