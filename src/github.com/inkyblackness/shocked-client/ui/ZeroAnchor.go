package ui

type zeroAnchor struct{}

// ZeroAnchor returns an anchor with the static value 0.
// Requests to set a new value are ignored.
func ZeroAnchor() Anchor {
	return &zeroAnchor{}
}

func (anchor *zeroAnchor) Value() float32 {
	return 0
}

func (anchor *zeroAnchor) RequestValue(newValue float32) {
}
