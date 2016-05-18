package video

import (
	check "gopkg.in/check.v1"
)

type PackControlWordsSuite struct {
}

var _ = check.Suite(&PackControlWordsSuite{})

func (suite *PackControlWordsSuite) TestEmptyArrayResultsInZeroCount(c *check.C) {
	data := PackControlWords(nil)

	c.Check(data, check.DeepEquals, []byte{0x00, 0x00, 0x00, 0x00})
}

func (suite *PackControlWordsSuite) TestCountIsInLengthTimesThree(c *check.C) {
	data := PackControlWords([]ControlWord{ControlWord(0), ControlWord(1)})

	c.Check(data[:4], check.DeepEquals, []byte{0x06, 0x00, 0x00, 0x00})
}

func (suite *PackControlWordsSuite) TestSingleItem(c *check.C) {
	data := PackControlWords([]ControlWord{ControlWord(0x00BBCCDD)})

	c.Check(data, check.DeepEquals, []byte{0x03, 0x00, 0x00, 0x00, 0xDD, 0xCC, 0xBB, 0x01})
}

func (suite *PackControlWordsSuite) TestMultiItemArray(c *check.C) {
	data := PackControlWords([]ControlWord{ControlWord(0x00BBCCDD), ControlWord(0x00112233)})

	c.Check(data, check.DeepEquals, []byte{0x06, 0x00, 0x00, 0x00, 0xDD, 0xCC, 0xBB, 0x01, 0x33, 0x22, 0x11, 0x01})
}

func (suite *PackControlWordsSuite) TestIdenticalWordsArePacked(c *check.C) {
	data := PackControlWords([]ControlWord{ControlWord(0x00BBCCDD), ControlWord(0x00BBCCDD)})

	c.Check(len(data), check.Equals, 8)
	c.Check(data[4:8], check.DeepEquals, []byte{0xDD, 0xCC, 0xBB, 0x02})
}

func (suite *PackControlWordsSuite) TestCountIsResetForFurtherWords(c *check.C) {
	data := PackControlWords([]ControlWord{ControlWord(0x00BBCCDD), ControlWord(0x00BBCCDD), ControlWord(0x00112233)})

	c.Check(len(data), check.Equals, 12)
	c.Check(data[4:12], check.DeepEquals, []byte{0xDD, 0xCC, 0xBB, 0x02, 0x33, 0x22, 0x11, 0x01})
}

func (suite *PackControlWordsSuite) TestMaximumCountIs255(c *check.C) {
	words := make([]ControlWord, 260)
	for i := 0; i < len(words); i++ {
		words[i] = ControlWord(0x00112233)
	}
	data := PackControlWords(words)

	c.Check(len(data), check.Equals, 12)
	c.Check(data, check.DeepEquals, []byte{0x0C, 0x03, 0x00, 0x00, 0x33, 0x22, 0x11, 0xFF, 0x33, 0x22, 0x11, 0x05})
}
