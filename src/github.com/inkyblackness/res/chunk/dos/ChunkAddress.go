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
	coder.CodeUint24(&addr.uncompressedLength)
	coder.CodeByte(&addr.chunkType)
	coder.CodeUint24(&addr.chunkLength)
	coder.CodeByte(&addr.contentType)
}
