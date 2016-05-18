package serial

import (
	check "gopkg.in/check.v1"
)

type PositioningEncoderSuite struct {
	coder PositioningCoder
	store *ByteStore
}

var _ = check.Suite(&PositioningEncoderSuite{})

func (suite *PositioningEncoderSuite) SetUpTest(c *check.C) {
	suite.store = NewByteStore()
	suite.coder = NewPositioningEncoder(suite.store)
}

func (suite *PositioningEncoderSuite) TestSetCurPosRepositionsWritePointer(c *check.C) {
	value := uint32(0)
	suite.coder.CodeUint32(&value)
	value = 0x13243546

	suite.coder.SetCurPos(0)
	suite.coder.CodeUint32(&value)
	result := suite.store.Data()

	c.Assert(result, check.DeepEquals, []byte{0x46, 0x35, 0x24, 0x13})
}
