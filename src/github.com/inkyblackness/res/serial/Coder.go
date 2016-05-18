package serial

// Coder represents an encoder/decoder for binary data
type Coder interface {
	// CodeUint16 serializes an unsigned 16bit integer value
	CodeUint16(value *uint16)
	// CodeUint24 serializes an unsigned 24bit integer value
	CodeUint24(value *uint32)
	// CodeUint32 serializes an unsigned 32bit integer value
	CodeUint32(value *uint32)
	// CodeBytes serializes a slice
	CodeBytes(value []byte)
	// CodeByte serializes a single byte
	CodeByte(value *byte)
}
