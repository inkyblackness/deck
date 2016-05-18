package core

type MaskingCoder struct {
	data   []byte
	offset int
}

func NewMaskingCoder(data []byte) *MaskingCoder {
	coder := &MaskingCoder{
		data:   data,
		offset: 0}

	return coder
}

func (coder *MaskingCoder) CodeUint16(value *uint16) {
	coder.mask(2)
}

func (coder *MaskingCoder) CodeUint24(value *uint32) {
	coder.mask(3)
}

func (coder *MaskingCoder) CodeUint32(value *uint32) {
	coder.mask(4)
}

func (coder *MaskingCoder) CodeBytes(value []byte) {
	coder.offset += len(value)
}

func (coder *MaskingCoder) CodeByte(value *byte) {
	coder.mask(1)
}

func (coder *MaskingCoder) mask(bytes int) {
	for endOffset := coder.offset + bytes; coder.offset < endOffset; coder.offset++ {
		coder.data[coder.offset] = 0x00
	}
}
