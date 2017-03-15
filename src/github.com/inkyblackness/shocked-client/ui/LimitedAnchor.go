package ui

type limitedAnchor struct {
	from      Anchor
	to        Anchor
	reference Anchor
}

// NewLimitedAnchor returns an anchor which limits the value of another
// anchor within two reference anchors.
// Requests to set a new value are forwarded to the reference anchor if the
// new value is within the allowed limits.
func NewLimitedAnchor(from, to, reference Anchor) Anchor {
	return &limitedAnchor{from: from, to: to, reference: reference}
}

func (anchor *limitedAnchor) Value() float32 {
	fromValue := anchor.from.Value()
	toValue := anchor.to.Value()
	referenceValue := anchor.reference.Value()
	result := referenceValue

	if toValue < fromValue {
		result = toValue + ((fromValue - toValue) / 2.0)
	} else if referenceValue < fromValue {
		result = fromValue
	} else if referenceValue > toValue {
		result = toValue
	}

	return result
}

func (anchor *limitedAnchor) RequestValue(newValue float32) {
	fromValue := anchor.from.Value()
	toValue := anchor.to.Value()

	if fromValue <= toValue {
		forwardedValue := newValue

		if forwardedValue < fromValue {
			forwardedValue = fromValue
		} else if forwardedValue > toValue {
			forwardedValue = toValue
		}
		anchor.reference.RequestValue(forwardedValue)
	}
}
