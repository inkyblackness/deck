package format

import (
	"bytes"
	"encoding/binary"

	check "gopkg.in/check.v1"
)

type HeaderSuite struct {
}

var _ = check.Suite(&HeaderSuite{})

func (suite *HeaderSuite) TestHeaderSerializesToProperLength(c *check.C) {
	source := bytes.NewReader(make([]byte, 0x200))
	var header Header

	binary.Read(source, binary.LittleEndian, &header)
	curPos, _ := source.Seek(0, 1)

	c.Check(curPos, check.Equals, int64(HeaderSize))
}
