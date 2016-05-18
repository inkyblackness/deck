package dos

import (
	"bytes"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"

	check "gopkg.in/check.v1"
)

type FormatReaderSuite struct {
}

var _ = check.Suite(&FormatReaderSuite{})

func (suite *FormatReaderSuite) TestNewProviderReturnsErrorOnNil(c *check.C) {
	_, err := NewProvider(nil, []objprop.ClassDescriptor{})

	c.Assert(err, check.ErrorMatches, "source is nil")
}

func (suite *FormatReaderSuite) TestNewProviderReturnsErrorOnFileWithWrongSize(c *check.C) {
	sourceData := append([]byte{0x2D, 0x00, 0x00, 0x00}, make([]byte, 10)...)

	_, err := NewProvider(bytes.NewReader(sourceData), []objprop.ClassDescriptor{})

	c.Assert(err, check.ErrorMatches, "Format mismatch")
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderOnEmptySource(c *check.C) {
	source := bytes.NewReader([]byte{0x2D, 0x00, 0x00, 0x00})
	provider, _ := NewProvider(source, []objprop.ClassDescriptor{})

	c.Assert(provider, check.NotNil)
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderForSampleData(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 1,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 2}}}
	class2 := objprop.ClassDescriptor{GenericDataLength: 2,
		Subclasses: []objprop.SubclassDescriptor{
			objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 1},
			objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 7}}}

	sourceData := getSamplePropertyData()

	source := bytes.NewReader(sourceData)
	provider, _ := NewProvider(source, []objprop.ClassDescriptor{class1, class2})

	c.Assert(provider, check.NotNil)
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderWithValidDataA(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 1,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 2}}}
	class2 := objprop.ClassDescriptor{GenericDataLength: 2,
		Subclasses: []objprop.SubclassDescriptor{
			objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 1},
			objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 7}}}

	sourceData := getSamplePropertyData()

	source := bytes.NewReader(sourceData)
	provider, _ := NewProvider(source, []objprop.ClassDescriptor{class1, class2})

	data := provider.Provide(res.MakeObjectID(res.ObjectClass(1), res.ObjectSubclass(0), res.ObjectType(0)))
	expected := objprop.ObjectData{
		Generic:  []byte{0x1A, 0x1B},
		Specific: []byte{0x10},
		Common:   make([]byte, objprop.CommonPropertiesLength)}
	expected.Common[0] = 0x33

	c.Assert(data, check.DeepEquals, expected)
}

func (suite *FormatReaderSuite) TestNewProviderReturnsProviderWithValidDataB(c *check.C) {
	class1 := objprop.ClassDescriptor{GenericDataLength: 1,
		Subclasses: []objprop.SubclassDescriptor{objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 2}}}
	class2 := objprop.ClassDescriptor{GenericDataLength: 2,
		Subclasses: []objprop.SubclassDescriptor{
			objprop.SubclassDescriptor{TypeCount: 1, SpecificDataLength: 1},
			objprop.SubclassDescriptor{TypeCount: 2, SpecificDataLength: 7}}}

	sourceData := getSamplePropertyData()

	source := bytes.NewReader(sourceData)
	provider, _ := NewProvider(source, []objprop.ClassDescriptor{class1, class2})

	data := provider.Provide(res.MakeObjectID(res.ObjectClass(1), res.ObjectSubclass(1), res.ObjectType(1)))
	expected := objprop.ObjectData{
		Generic:  []byte{0x3A, 0x3B},
		Specific: []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36},
		Common:   make([]byte, objprop.CommonPropertiesLength)}
	expected.Common[0] = 0x55

	c.Assert(data, check.DeepEquals, expected)
}

func getSamplePropertyData() []byte {
	sourceData := []byte{0x2D, 0x00, 0x00, 0x00} // header
	// class 0
	sourceData = append(sourceData, 0xAA, 0xBA)             // 2x class 0 type generic data
	sourceData = append(sourceData, 0xA0, 0xA1, 0xB0, 0xB1) // 2x class 0|0 types 0 and 1 specific data
	// class 1
	sourceData = append(sourceData, 0x1A, 0x1B, 0x2A, 0x2B, 0x3A, 0x3B)       // 3x class 1 type generic data
	sourceData = append(sourceData, 0x10)                                     // 1x class 1|0 type 0 specific data
	sourceData = append(sourceData, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26) // class 1|1 type 0 specific data
	sourceData = append(sourceData, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36) // class 1|1 type 1 specific data
	// 5x common data
	for i := 1; i <= 5; i++ {
		sourceData = append(sourceData, byte(i<<4|i))
		sourceData = append(sourceData, make([]byte, objprop.CommonPropertiesLength-1)...)
	}

	return sourceData
}
