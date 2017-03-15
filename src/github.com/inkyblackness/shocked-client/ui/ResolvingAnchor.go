package ui

// AnchorResolver is a function returning another anchor.
type AnchorResolver func() Anchor

type resolvingAnchor struct {
	resolver AnchorResolver
}

// NewResolvingAnchor returns an anchor resolving to another anchor via a resolver function.
// Requests to set a new value are forwarded.
func NewResolvingAnchor(resolver AnchorResolver) Anchor {
	return &resolvingAnchor{resolver: resolver}
}

func (anchor *resolvingAnchor) resolved() Anchor {
	resolved := anchor.resolver()
	if resolved == nil {
		resolved = ZeroAnchor()
	}
	return resolved
}

func (anchor *resolvingAnchor) Value() float32 {
	return anchor.resolved().Value()
}

func (anchor *resolvingAnchor) RequestValue(newValue float32) {
	anchor.resolved().RequestValue(newValue)
}
