package ui

type offsetAnchor struct {
	base   Anchor
	offset float32
}

// NewOffsetAnchor returns an anchor with an absolute offset to another
// anchor.
// Requests to set a new value change the offset.
func NewOffsetAnchor(base Anchor, offset float32) Anchor {
	return &offsetAnchor{base: base, offset: offset}
}

func (anchor *offsetAnchor) Value() float32 {
	return anchor.base.Value() + anchor.offset
}

func (anchor *offsetAnchor) RequestValue(newValue float32) {
	anchor.offset = newValue - anchor.base.Value()
}
