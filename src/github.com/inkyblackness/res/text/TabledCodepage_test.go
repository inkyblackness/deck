package text

import (
	check "gopkg.in/check.v1"
)

type TabledCodepageSuite struct {
	cp Codepage
}

var _ = check.Suite(&TabledCodepageSuite{})

func (suite *TabledCodepageSuite) SetUpTest(c *check.C) {
	suite.cp = DefaultCodepage()
}

func (suite *TabledCodepageSuite) TestEncode(c *check.C) {
	result := suite.cp.Encode("ä")

	c.Check(result, check.DeepEquals, []byte{132, 0x00})
}

func (suite *TabledCodepageSuite) TestDecode(c *check.C) {
	result := suite.cp.Decode([]byte{212, 225, 0x00})

	c.Check(result, check.Equals, "Èß")
}
