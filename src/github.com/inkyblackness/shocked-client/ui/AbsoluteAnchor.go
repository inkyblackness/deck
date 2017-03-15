package ui

type absoluteAnchor struct {
	value float32
}

// NewAbsoluteAnchor returns an anchor with a determined value.
// Requests to set a new value are directly applied.
func NewAbsoluteAnchor(value float32) Anchor {
	return &absoluteAnchor{value: value}
}

func (anchor *absoluteAnchor) Value() float32 {
	return anchor.value
}

func (anchor *absoluteAnchor) RequestValue(newValue float32) {
	anchor.value = newValue
}
