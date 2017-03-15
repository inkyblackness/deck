package ui

type relativeAnchor struct {
	from     Anchor
	to       Anchor
	fraction float32
}

// NewRelativeAnchor returns an anchor which derives a value from a fraction
// between two other anchors.
// Requests to set a new value update the fraction value.
func NewRelativeAnchor(from, to Anchor, fraction float32) Anchor {
	return &relativeAnchor{from: from, to: to, fraction: fraction}
}

func (anchor *relativeAnchor) Value() float32 {
	fromValue := anchor.from.Value()

	return fromValue + (anchor.to.Value()-fromValue)*anchor.fraction
}

func (anchor *relativeAnchor) RequestValue(newValue float32) {
	fromValue := anchor.from.Value()
	toValue := anchor.to.Value()

	if fromValue != toValue {
		anchor.fraction = (newValue - fromValue) / (toValue - fromValue)
	}
}
