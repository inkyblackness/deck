package dos

import (
	"github.com/inkyblackness/res/serial"
)

type chunkAddress struct {
	startOffset uint32

	chunkLength        uint32
	uncompressedLength uint32
	chunkType          byte
	contentType        byte
}

func (addr *chunkAddress) code(coder serial.Coder) {
	fieldA := (addr.uncompressedLength & 0x00FFFFFF) | uint32(addr.chunkType)<<24
	coder.Code(&fieldA)
	addr.uncompressedLength = fieldA & 0x00FFFFFF
	addr.chunkType = byte(fieldA >> 24)
	fieldB := (addr.chunkLength & 0x00FFFFFF) | uint32(addr.contentType)<<24
	coder.Code(&fieldB)
	addr.chunkLength = fieldB & 0x00FFFFFF
	addr.contentType = byte(fieldB >> 24)
}
