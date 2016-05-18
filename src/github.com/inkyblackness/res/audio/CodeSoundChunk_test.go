package audio

import (
	"github.com/inkyblackness/res/audio/mem"

	check "gopkg.in/check.v1"
)

type CodeSoundChunkSuite struct {
}

var _ = check.Suite(&CodeSoundChunkSuite{})

func (suite *CodeSoundChunkSuite) TestChunkTypeReturnsProvidedValue(c *check.C) {
	rawSoundSamples := []byte{0x00, 0x20, 0x40, 0x80, 0xC0, 0xFF}
	rawSound := mem.NewL8SoundData(20000.0, rawSoundSamples)
	encoded := EncodeSoundChunk(rawSound)
	decoded, err := DecodeSoundChunk(encoded)

	c.Assert(err, check.IsNil)

	c.Check(decoded.SampleRate(), check.Equals, rawSound.SampleRate())
	c.Check(decoded.SampleCount(), check.Equals, rawSound.SampleCount())
	samples := decoded.Samples(0, decoded.SampleCount())
	c.Check(samples, check.DeepEquals, rawSoundSamples)
}
