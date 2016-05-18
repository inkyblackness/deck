package dos

import (
	"bytes"

	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/textprop"

	check "gopkg.in/check.v1"
)

type FormatWriterSuite struct {
}

var _ = check.Suite(&FormatWriterSuite{})

func (suite *FormatWriterSuite) TestWriterCreatesCompatibleData(c *check.C) {
	prop0 := createTestProperties(0x01)
	prop1 := createTestProperties(0x02)
	prop2 := createTestProperties(0x03)
	store := serial.NewByteStore()
	consumer := NewConsumer(store)

	consumer.Consume(0, prop0)
	consumer.Consume(1, prop1)
	consumer.Consume(2, prop2)
	consumer.Finish()

	provider, _ := NewProvider(bytes.NewReader(store.Data()))

	c.Assert(provider.Provide(1), check.DeepEquals, prop1)
}

func (suite *FormatWriterSuite) TestWriterCreatesScratchBuffer(c *check.C) {
	store := serial.NewByteStore()
	consumer := NewConsumer(store)

	consumer.Finish()
	provider, _ := NewProvider(bytes.NewReader(store.Data()))

	c.Assert(provider.EntryCount(), check.DeepEquals, uint32(363))
}

func createTestProperties(filler byte) []byte {
	data := make([]byte, textprop.TexturePropertiesLength)

	for i := 0; i < len(data); i++ {
		data[i] = filler
	}

	return data
}
