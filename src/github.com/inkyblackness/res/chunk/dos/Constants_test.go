package dos

import (
	"github.com/inkyblackness/res/serial"
)

func emptyResourceFile() []byte {
	store := serial.NewByteStore()
	encoder := serial.NewPositioningEncoder(store)

	codeHeader(encoder)
	// write offset to dictionary - in this case right after header
	{
		dictionaryOffset := uint32(store.Len() + 4)
		encoder.CodeUint32(&dictionaryOffset)
	}
	{
		numberOfChunks := uint16(0)
		firstChunkOffset := uint32(store.Len())

		encoder.CodeUint16(&numberOfChunks)
		encoder.CodeUint32(&firstChunkOffset)
	}

	return store.Data()
}
