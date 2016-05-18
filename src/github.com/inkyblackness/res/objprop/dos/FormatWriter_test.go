package dos

import (
	"bytes"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/serial"

	check "gopkg.in/check.v1"
)

type FormatWriterSuite struct {
}

var _ = check.Suite(&FormatWriterSuite{})

func (suite *FormatWriterSuite) TestWriterCreatesZeroDataWhenConsumeIsMissing(c *check.C) {
	descriptor := []objprop.ClassDescriptor{}
	store := serial.NewByteStore()
	consumer := NewConsumer(store, descriptor)

	consumer.Finish()
	provider, _ := NewProvider(bytes.NewReader(store.Data()), descriptor)

	c.Assert(provider, check.NotNil)
}

func (suite *FormatWriterSuite) TestWriterStoresDataAccordingToFormat(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 3,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 2}}}
	class2 := objprop.ClassDescriptor{GenericDataLength: 5,
		Subclasses: []objprop.SubclassDescriptor{
			objprop.SubclassDescriptor{TypeCount: 5, SpecificDataLength: 13},
			objprop.SubclassDescriptor{TypeCount: 3, SpecificDataLength: 4}}}
	descriptor := []objprop.ClassDescriptor{class1, class2}
	store := serial.NewByteStore()
	consumer := NewConsumer(store, descriptor)

	id := res.MakeObjectID(1, 0, 4)
	expected := objprop.ObjectData{
		Generic:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		Specific: []byte{0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD},
		Common:   make([]byte, objprop.CommonPropertiesLength)}

	consumer.Consume(id, expected)

	consumer.Finish()
	provider, _ := NewProvider(bytes.NewReader(store.Data()), descriptor)
	retrieved := provider.Provide(id)

	c.Assert(retrieved, check.DeepEquals, expected)
}

func (suite *FormatWriterSuite) TestConsumePanicsIfGenericLengthMismatch(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 3,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 2}}}
	descriptor := []objprop.ClassDescriptor{class1}
	store := serial.NewByteStore()
	consumer := NewConsumer(store, descriptor)

	id := res.MakeObjectID(0, 0, 0)
	expected := objprop.ObjectData{
		Generic:  []byte{0x01, 0x02, 0x03, 0x04},
		Specific: []byte{0xA1, 0xA2},
		Common:   make([]byte, objprop.CommonPropertiesLength)}

	c.Assert(func() { consumer.Consume(id, expected) }, check.Panics, errSizeMismatch)
}

func (suite *FormatWriterSuite) TestConsumePanicsIfSpecificLengthMismatch(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 3,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 2}}}
	descriptor := []objprop.ClassDescriptor{class1}
	store := serial.NewByteStore()
	consumer := NewConsumer(store, descriptor)

	id := res.MakeObjectID(0, 0, 0)
	expected := objprop.ObjectData{
		Generic:  []byte{0x01, 0x02, 0x03},
		Specific: []byte{0xA1},
		Common:   make([]byte, objprop.CommonPropertiesLength)}

	c.Assert(func() { consumer.Consume(id, expected) }, check.Panics, errSizeMismatch)
}

func (suite *FormatWriterSuite) TestConsumePanicsIfCommonLengthMismatch(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 3,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 2}}}
	descriptor := []objprop.ClassDescriptor{class1}
	store := serial.NewByteStore()
	consumer := NewConsumer(store, descriptor)

	id := res.MakeObjectID(0, 0, 0)
	expected := objprop.ObjectData{
		Generic:  []byte{0x01, 0x02, 0x03},
		Specific: []byte{0xA1, 0xA2},
		Common:   make([]byte, objprop.CommonPropertiesLength+1)}

	c.Assert(func() { consumer.Consume(id, expected) }, check.Panics, errSizeMismatch)
}
