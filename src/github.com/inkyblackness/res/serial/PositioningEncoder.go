package serial

import "io"

type positioningEncoder struct {
	encoder

	seeker io.Seeker
}

// NewPositioningEncoder returns an encoder that also implements the Positioner interface.
func NewPositioningEncoder(dest io.WriteSeeker) PositioningCoder {
	coder := &positioningEncoder{encoder: encoder{dest: dest, offset: 0}, seeker: dest}

	return coder
}

// CurPos gets the current position in the data
func (coder *positioningEncoder) CurPos() uint32 {
	return coder.offset
}

// SetCurPos sets the current position in the data
func (coder *positioningEncoder) SetCurPos(offset uint32) {
	coder.seeker.Seek(int64(offset), 0)
	coder.offset = offset
}
