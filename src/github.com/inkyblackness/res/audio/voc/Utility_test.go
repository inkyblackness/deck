package voc

import (
	check "gopkg.in/check.v1"
)

type UtilitySuite struct {
}

var _ = check.Suite(&UtilitySuite{})

func (suite *UtilitySuite) TestFrequencyDivisorFor11111(c *check.C) {
	divisor := byte(0xA6)
	sampleRate := divisorToSampleRate(divisor)
	result := sampleRateToDivisor(sampleRate)

	c.Check(result, check.Equals, divisor)
}

func (suite *UtilitySuite) TestFrequencyDivisorFor222222(c *check.C) {
	divisor := byte(0x2D)
	sampleRate := divisorToSampleRate(divisor)
	result := sampleRateToDivisor(sampleRate)

	c.Check(result, check.Equals, divisor)
}
