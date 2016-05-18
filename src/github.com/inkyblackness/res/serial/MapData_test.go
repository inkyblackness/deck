package serial

import (
	"bytes"

	check "gopkg.in/check.v1"
)

type TestData struct {
	SomeUint8 byte
	SomeInt8  int8

	SomeUint16 uint16
	SomeInt16  int16

	SomeUint32 uint32
	SomeInt32  int32

	SomeByteArray [5]byte
}

type TestDataList struct {
	Entries []*TestData
}

type TestDataString struct {
	SingleEntry string
}

type TestDataNested struct {
	TestData

	ExtraByteArray [3]byte
	SomeInts       [4]int16
}

type MapDataSuite struct {
}

var _ = check.Suite(&MapDataSuite{})

func (suite *MapDataSuite) SetUpTest(c *check.C) {
}

func (suite *MapDataSuite) verifyMapData(c *check.C, v interface{}, data []byte) {
	store := NewByteStore()

	MapData(v, NewDecoder(bytes.NewReader(data)))
	MapData(v, NewEncoder(store))
	c.Assert(store.Data(), check.DeepEquals, data)
}

func (suite *MapDataSuite) TestMapDataCodesZeroValues(c *check.C) {
	v := &TestData{}
	data := make([]byte, 19)

	suite.verifyMapData(c, v, data)
}

func (suite *MapDataSuite) TestMapDataCodesNegativeValues(c *check.C) {
	v := &TestData{}
	data := make([]byte, 19)

	for i := range data {
		data[i] = 0x80
	}

	suite.verifyMapData(c, v, data)
}

func (suite *MapDataSuite) TestMapDataCodesAllBitsSet(c *check.C) {
	v := &TestData{}
	data := make([]byte, 19)

	for i := range data {
		data[i] = 0xFF
	}

	suite.verifyMapData(c, v, data)
}

func (suite *MapDataSuite) TestMapDataCodesArrayValues(c *check.C) {
	v := &TestDataList{Entries: []*TestData{&TestData{}, &TestData{}}}
	data := make([]byte, 19*len(v.Entries))

	suite.verifyMapData(c, v, data)
}

func (suite *MapDataSuite) TestMapDataCodesStringValues(c *check.C) {
	v := &TestDataString{}
	data := []byte{0x31, 0x32, 0x33, 0x34, 0x00}

	suite.verifyMapData(c, v, data)
}

func (suite *MapDataSuite) TestMapDataCodesNestedValues(c *check.C) {
	v := &TestDataNested{}
	data := make([]byte, 19+3+4*2)

	suite.verifyMapData(c, v, data)
}
