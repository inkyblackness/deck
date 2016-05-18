package serial

import "io"

// encoder implements the Coder interface to write to a writer
type encoder struct {
	offset uint32
	dest   io.Writer
}

// NewEncoder creates and returns a fresh encoder
func NewEncoder(dest io.Writer) Coder {
	coder := &encoder{
		offset: 0,
		dest:   dest}

	return coder
}

// CodeByte encodes a single byte
func (coder *encoder) CodeByte(value *byte) {
	coder.writeBytes(*value)
}

// CodeBytes encodes the provided bytes
func (coder *encoder) CodeBytes(value []byte) {
	coder.writeBytes(value...)
}

// CodeUint16 encodes an unsigned 16bit value
func (coder *encoder) CodeUint16(value *uint16) {
	coder.writeBytes(byte((*value>>0)&0xFF), byte((*value>>8)&0xFF))
}

// CodeUint24 encodes an unsigned 24bit value
func (coder *encoder) CodeUint24(value *uint32) {
	coder.writeBytes(byte((*value>>0)&0xFF), byte((*value>>8)&0xFF), byte((*value>>16)&0xFF))
}

// CodeUint32 encodes an unsigned 32bit value
func (coder *encoder) CodeUint32(value *uint32) {
	coder.writeBytes(byte((*value>>0)&0xFF), byte((*value>>8)&0xFF), byte((*value>>16)&0xFF), byte((*value>>24)&0xFF))
}

func (coder *encoder) writeBytes(bytes ...byte) {
	written, err := coder.dest.Write(bytes)
	coder.offset += uint32(written)
	if err != nil {
		panic(err)
	}
}
