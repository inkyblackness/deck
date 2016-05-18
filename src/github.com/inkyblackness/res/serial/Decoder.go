package serial

import "io"

// decoder is for decoding from a random access stream
type decoder struct {
	source io.Reader
	offset uint32
}

// NewDecoder creates a new decoder from given source
func NewDecoder(source io.Reader) Coder {
	coder := &decoder{source: source, offset: 0}

	return coder
}

// CodeByte decodes a single byte
func (coder *decoder) CodeByte(value *byte) {
	buf := coder.readBytes(1)
	*value = buf[0]
}

// CodeBytes decodes the provided bytes
func (coder *decoder) CodeBytes(value []byte) {
	read, err := coder.source.Read(value)
	coder.offset += uint32(read)
	if err != nil {
		panic(err)
	}
}

// CodeUint16 decodes an unsigned 16bit value
func (coder *decoder) CodeUint16(value *uint16) {
	buf := coder.readBytes(2)
	*value = (uint16(buf[0]) << 0) | (uint16(buf[1]) << 8)
}

// CodeUint24 decodes a 24bit unsigned integer
func (coder *decoder) CodeUint24(value *uint32) {
	buf := coder.readBytes(3)
	*value = (uint32(buf[0]) << 0) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16)
}

// CodeUint32 decodes a 32bit unsigned integer
func (coder *decoder) CodeUint32(value *uint32) {
	buf := coder.readBytes(4)
	*value = (uint32(buf[0]) << 0) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16) | (uint32(buf[3]) << 24)
}

func (coder *decoder) readBytes(size int) []byte {
	buf := make([]byte, size)
	coder.CodeBytes(buf)

	return buf
}
