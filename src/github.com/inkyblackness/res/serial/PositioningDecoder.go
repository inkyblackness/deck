package serial

import "io"

// positioningDecoder is for decoding from a random access stream
type positioningDecoder struct {
	decoder

	seeker io.Seeker
}

// NewPositioningDecoder creates a new decoder from given source
func NewPositioningDecoder(source io.ReadSeeker) PositioningCoder {
	coder := &positioningDecoder{decoder: decoder{source: source, offset: 0}, seeker: source}

	return coder
}

// CurPos gets the current position in the data
func (coder *positioningDecoder) CurPos() uint32 {
	return coder.offset
}

// SetCurPos sets the current position in the data
func (coder *positioningDecoder) SetCurPos(offset uint32) {
	coder.seeker.Seek(int64(offset), 0)
	coder.offset = offset
}
