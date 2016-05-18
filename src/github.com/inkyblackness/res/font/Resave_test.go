package font

import (
	"bytes"
	"encoding/binary"

	check "gopkg.in/check.v1"
)

type ResaveSuite struct {
}

var _ = check.Suite(&ResaveSuite{})

func (suite *ResaveSuite) TestLoadSaveLoadSaveCreatesSameData(c *check.C) {
	data1 := suite.aSimpleFont()
	newFont, err := Load(bytes.NewReader(data1))

	c.Assert(err, check.IsNil)

	data2 := Save(newFont)

	c.Check(data2, check.DeepEquals, data1)
}

func (suite *ResaveSuite) aSimpleFont() []byte {
	var header Header
	buf := bytes.NewBuffer(nil)

	header.Type = Color
	header.FirstCharacter = 32
	header.LastCharacter = 32
	header.XOffsetStart = uint32(HeaderSize)
	header.BitmapStart = header.XOffsetStart + uint32(4)
	header.Width = 1
	header.Height = 1
	binary.Write(buf, binary.LittleEndian, &header)
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	buf.Write([]byte{0xAB})

	return buf.Bytes()
}
